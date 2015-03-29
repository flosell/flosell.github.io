---
layout: post
title:  "Introducing LambdaCD"
categories: lambdacd
date:   2015-03-14 13:26:16
---

The problem
===========

Do you have a build pipeline that is more complex than just `build`, `test`, `deploy`?
Are you happy with it? Do you fully understand what all those build-steps are doing, what input they take and where this input comes from? Do you like all those fragile plugins you need just to convince your CI server to do the thing you want to do? Do you like clicking your way through a web interface to get stuff done? Is your build pipeline versioned? Is it _tested_?

Or maybe your organization is trying out microservices. How easy is it for you to set up a pipeline for a new microservice? Have you figured out how to reuse aspects that are common for all services, how to share them between teams and still have stable interfaces, abstractions and dependencies?


An idea
=======

As developers, we know those problems are already solved in other parts of our world: Our code is readable, it has abstractions, testing, dependency management and versioning works effortlessly. Or servers are immutable and their deployment automated, everything built with the same concepts and tools as is our code. 

One thing is missing: The glue, the build pipeline between a commit and your critical live systems. The thing that, supposedly, [is the highest priority][fowler-ci-fix-immediately] of the development team. 

What we need is a way to put our build pipelines where our business logic, our database-layout and our server configuration already is: **in code**!

What does that get us? It gets us everything we love about code: We can structure it, version it, test it, share it with others any way we like! We can take advantage of all the tools and libraries we know and love!

You want to find out if your deployment really stops when a smoketest fails on staging? A small test with the right mock will tell you!

All your development teams should use the same process to deploy to your private cloud? Just put everything into a small library for everyone to include in their pipelines!

Your management requires your pipeline to behave differently on mondays, except for rainy days? You can build that! (not saying you should)

Simply put: Our build pipeline can be just like any other piece of software we release, with the same tools, methods and quality we are used to! 


[fowler-ci-fix-immediately]: http://www.martinfowler.com/articles/continuousIntegration.html#FixBrokenBuildsImmediately