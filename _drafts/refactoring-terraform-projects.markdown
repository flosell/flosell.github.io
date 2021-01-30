---
layout: post
title:  "Refactoring Terraform Code"
categories: terraform infrastructure-as-code refactoring

---

As developers, we do it without even thinking about it: We change things, rename, restructure our code without even thinking about it twice. It keeps our code clean. It's an acknowledgement that the that the code we wrote two months ago might no longer be fit for purpose now. That we learned new things since then. That our project grew and and the requirements on our codebase with them.

So if we are serious about saying that infrastructure is code, we need to do the same things. Unfortunately, as always, things are a bit more tricky in the infrastructure world. And even more so in the production infrastructure world where you can't simply revert your commit, tear down the world and start over. 

So in this post I will cover different strategies I have used to restructure terraform code without getting your state all messed up. I will talk about renaming resources, moving them into or between modules and merging or splitting whole projects. 
I will, for the purpose of this post, ignore the question how you'll end up with a well-structured terraform codebase. That'll have to wait for another post.

# The beginning

So for now, let's just assume you cleaned up your code a bit, ran `terraform plan` and are now stuck with the dreaded wall of red and green, followed by

```
Plan: 34 to add, 0 to change, 34 to destroy. 
``` 

This was probably not what you intended. You wanted to make your code easier to understand and maintain, not rebuild your whole infrastructure. So from here on out, our goal will be simple: Get back to an empty `terraform plan`. 

Small aside (TODO: format, structure somehow differently?): I'm assuming here that your code changes were **"safe refactorings"** like moves and renames. Those don't change the behaviour of our code so we would expect to see no impact on our terraform plan or our tests (TODO: really the definition?). If you want to make changes that will change behaviour, do them separately, before or after your safe refactoring steps. Otherwise, you'll just make this so much harder (TODO: more detail?)

With that out of the way, let's get started! 

# Approach 1: Just let terraform do its job

This is going to be the most straightforward and, depending on your circumstances, the safest approach: Just apply the plan, let terraform figure it out!

The pros of this approach: You have done this a thousand times, you know how this works. There's no special magic, brainpower or digging in terraforms internals involved. Terraform will do exactly what it told you it would do. So if the plan looks good to you, just apply it! 

The cons: If you googled hard enough to find this blog post, chances are there's something in the plan that you don't like. Maybe it would try to destroy your production database or regenerate all your IAM users passwords in the process. Maybe your resources depend on each other in a way that makes your plan had to apply the way it is. 

All of these things happen so we need ways to coax terraform back into a normal state. 

# Disclaimer: Danger ahead

(TODO: format differently?)

This is a warning. Approach 1 is the only approach in here that follows the standard terraform methods that you use in your day-to-day life. From here on out, you are in fully manual flight using advanced features, untested tools and error-prone techniques. In short, **you'll get chances to screw things up.** So make sure you come prepared:

Pull the latest state of the source code. Have a backup of every terraform state you are planning to touch (or enable versioning if you store them in S3). Get yourself a cup of your favorite (non-alcoholic) beverage, make sure you are focused. Read commands and outputs carefully, know what's going to happen.

If you normally apply changes through automation (as you should!), now is a good time to pause it. If you have multiple people working on your infrastructure, tell them to leave the pieces you'll touch alone for a while. You don't want other changes overwriting yours or picking up an intermediate state.

Ready? Are you _really_ sure you have your backups ready? Let's go! 

# Approach 2: Moving/renaming resource using `terraform state mv`

If you just renamed resources or moved them from one module to another within the same project, `terraform state mv` is the command to go for. It's the built-in way of telling terraform: "The resource you know under this name now has a different name". So if you moved an IAM user from the root of your project into a module, this will look something like this:
 
```
$ terraform state mv aws_iam_user.some_user module.human_users.aws_iam_user.some_user
```

If you want to get an idea which resources you'll need to move, a little bash can give you an idea: 

```
$ terraform plan -no-color | grep '^[[:space:]]*[+-]'
```

For a small number of resources, this is a decent, relatively safe way of getting your state back in sync with your code. 
However, if you moved around a larger part of your code base (e.g. split up a larger project into many modules), this can become tedious (and therefore error-prone) really quickly and you are probably looking for a way to generate those move-commands for you. 

One approach that can work is to expand on the bash command above to create a little script to do that for you. 
This works well enough but can be a little tricky to get right (especially if you aren't a scripting ninja).

Thankfully, the open-source community already did some of that work for you: [`tfmv`](https://github.com/afeld/tfmv) does pretty much that.
(TODO:write)
(TODO: link to https://ryaneschinger.com/blog/terraform-state-move/)

(TODO: extension with `tfmv`)

# Approach 3: Re-importing an existing resource using `terraform import`

Sometimes, a resource exists in your infrastructure but not in your terraform code. This might be the case when you are migrating manually configured infrastructure into terraform or because you moved a resource from one project to another.

In this case, `terraform import` allows you to import the state for a resource described in the terraform code. In short, you'll be mapping an identifier used inside your infrastructure platform (e.g. an EC2 Instance-ID) to the address of the resource in terraform:

```
$ terraform import aws_iam_policy.datadog DatadogAWSIntegrationPolicy
```

The details will be slightly different for each resource (with some not supporting the feature at all) so make sure you check your providers' documentation. For example, Route53 records are identified by the zone id concatenated with the FQDN: (TODO: check wording! fqdn? zone-id?) 

```
$ terraform import aws_route53_record.some_record ABCDEF123T7X3Y4_example.com
```

After you have imported a resource, it's a good idea to run `terraform plan` to make sure everything was imported correctly and your code matches what was imported. Ideally, you'd not see a diff at all. However, imports aren't perfect: Sometimes, a diff might remain simply because not everything. Most of the time, the next `apply` will sort this out (however, as always, read the plan carefully before applying to avoid nasty surprises).

Optional: If you imported a resource that already exists in a different project (e.g. because you moved it between projects), make sure you remove it from the old project and use `terraform state rm` to remove it from the old state. The same resource being tracked by two different projects is a recipe for chaos!  

(TODO:write under the title "moving individual resources between projects" (or add a subsection "case study, moving resources between projects"))
(TODO: extension with `terraforming`)
(TODO: maybe warn about associations that can't easily be imported?)

# Approach 4: Moving and merging whole state files using `terraform state pull/push`

When attempting bigger refactorings (e.g. when merging two projects into one or doing major restructurings), operations on individual resources are too tedious or just don't make sense. In this case, it can make sense to work on the complete statefile: `terraform state pull` will print the complete state in JSON-format to stdout. `terraform state push` will take such a state file and push it back into it's backing store in a consistent state. Notice that most of the time, terraform will consider such a push dangerous (which it is, you are overwriting the known state with arbitrary data). Use `-force` if you need to override this warning.

(TODO:write more (why not directly in backing storage? editing file for merging; be careful, this is internal state, it might change between terraform versions, ...))
(Note: useful for example when switching between workspaces and module based structure (maybe just add this as a case study section?))
 
# Approach 5: Manually editing state-files

(TODO:write)
(Note: this is an extension of the push/pull)
(Note: Do use push/pull, do not edit files directly in the backing store (e.g. S3), you might run into inconsistency issues (e.g. with dynamodb))
(Note: only do this if nothing else helps)


# Postscript (TODO: right naming?): How to recover once you screwed up

TODO
