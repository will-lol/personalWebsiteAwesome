---
title: You don't need the edge
date: 2024-02-24
slug: edge
description: >
  Edge computing makes the concept of \'the cloud\' real. But is it actually that different from what we had before?
---
In 2006, AWS launched three products. Two of them changed computing. S3 and EC2 introduced the idea of 'the cloud' to developers. A concept that elucidates ideas of magical, floating computers that were everywhere you needed them to be. SQS was also released-someone might argue its equal importance. 

Until the popularisation of edge computing by Vercel, computers remained physically planted. It was obvious to both developers and Australians that the computers they interacted with in (mostly) US-East-1 were definitely not omnipresent. While possible before, edge computing made the idea of running arbitrary code on a server close to your users very easy-and fast. I think this magical experience realised the concept of 'the cloud', but do you really need it for your side project of 3 users?

No. In fact, 'the edge' is rather limited in its use case. 

'The edge' is currently being used primarily for the server side rendering of websites. I will examine this use case and determine if I think the benefits of edge computing warrant its (quite high) price tag. At the time of writing, Lambda functions cost USD$0.2 per 1 million invocations and Lambda@Edge costs a whopping USD$0.6 per 1 million invocations.

As opposed to a static approach, server side rendering provides one main benefit: dynamic content insertion in the files themselves. This might reduce the amount of time it takes for a user to see a fully rendered page. Running your servers 'on the edge' will provide even greater benefits.

A static approach however will always show a visual result (if incomplete) before server side rendering can. If data fetching is well handled, then API requests to our slow(er) serverless endpoints shouldn't be so detrimental. After all, the user can at least see a placeholder or loading bar.

Where the dynamic data is not critical, traditional serverless computing could be combined with a caching layer to provide outdated information that could be updated on the client. This means the benefit of edge is made so small that its costs may be unjustifiable. 

Edge computing also provides a number of limitations that make it unusable in some circumstances. Most implementations are limited to JavaScript runtimes. In **almost all** cases, developers will find themselves making requests to other endpoints from the edge runtime. If you *could* have made the request in a JavaScript runtime, why make the client wait for your server when the it can show indications of progress if it handles the data fetching. Many of these endpoints may not be running on the edge, slowing the initial page load even further.

I might also question *why* we are chasing the fastest page load times. Like everyone, I would love for the web to feel instant, but when [Walmart measures 2% more conversions for every **1 second** of page load improvement](https://wwww.cloudflare.com/en-au/learning/performance/more/website-performance-conversion-rates/), I wonder if the 100ms improvement you may see from switching to edge *really* justifies its price.

I do not recommend that you move your database SDK to the client, but hitting a serverless endpoint on the client instead of in an edge function at page load time may result in a perceived speed *improvement*-especially if your database isn't replicated in your user's region. 

I think that the limitations of 'the edge' reduce its use case so severely that I would seldom use it at all. Unless you are building a speed critical CRUD dashboard with an edge database. 
