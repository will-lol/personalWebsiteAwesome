---
title: CSS Margin, Gap, and Flexbox
date: 2024-04-25
slug: gap
description: >
  CSS margins are a major source of headache for web developers. Is there a solution?
---
**I like the CSS box model.**

![box model representation](/assets/images/box-model.svg)

It makes sense. There’s only one issue: that's not it.

The CSS box model actually has the borders on the _outside_ of the box. So when you add a border, it will overflow the element! 

![diagram showing how borders extend outside of the parent](/assets/images/border-box.svg)

As developers, we get the `box-sizing: border-box` rule to solve this issue.

We rarely come across the case where the borders of an element are so large that they cause our layouts to break, so `box-sizing: border-box` (despite being a more sensible default) is often [forgotten](https://css-tricks.com/international-box-sizing-awareness-day/).

However, there is one other (much more common) property that is so ubiquitous that its odd side effects are unquestioned. The CSS margin property’s tendency to push sibling elements around is a squabble that we seem to largely put up with as developers.

Margins are most often used to add arbitrary amounts of space between elements to make them _feel_ right together. There is often no rhyme or reason to the margins we add to elements.

The idea that there is a CSS property that we arbitrarily tinker with until the page looks correct seems fundamentally at odds with the way we architect web applications. Inevitably, margins worm their way into our classes and components. This makes them less reusable. An element’s margin should be a property of the particular composition of elements, _not_ the element itself.


## What if we didn’t use margins

Perhaps the most simple solution to margin woes is to never use them. One must resist the urge to replace their margins with padding, as padding has a much different purpose and _makes sense to be _a property of the element. A property much more fit for purpose (and promising!) is `gap`.

`gap` is a property set on a `grid` or `flex` parent. It spaces elements out evenly with a `gap` between them. I think CSS `gap` might be my favourite property. I love it. 

Unlike margins, `gap` _is_ a property of a group of elements. This makes it very well suited to `margin`’s job—and I think it should almost always replace margins in situations where a uniform gap is needed between elements. However, as much as I wish I could set `gap` on elements using the `block` layout algorithm, the added noise of `display: flex | grid` can muddy up CSS classes a bit—but it's a price I’m willing to pay. 

`gap` is great until you come up with the following bright idea:

```html
<article style="display: flex;flex-direction: column;gap: 1rem;">
  <h1>Hello my friends</h1>
  <p>This is a call to all my friends</p>
  <h2>I have a lot of them</h2>
  <p>I do</p>
  <h2>Are you free on wednesday I want to do a sleepover</h2>
  <p> ok I will come</p>
</article>
```

This particular example looks… ok—but designers will cringe at this idea.

Elements in typography should _not_ have uniform gaps between them. It's important for readability that headings and paragraphs are less spaced apart than paragraphs and paragraphs! We’ve also just added `flex` to this sequence of elements that is clearly more suited to the good old fashioned `block` layout!

It's not just typography either—there are quite a few situations where we might want uneven gaps between elements.


## `<article>`, `display: block`, and margins

From MDN:

> The `<article>` [HTML](https://developer.mozilla.org/en-US/docs/Web/HTML) element represents a self-contained composition in a document, page, application, or site, which is intended to be independently distributable or reusable (e.g., in syndication). Examples include: a forum post, a magazine or newspaper article, or a blog entry, a product card, a user-submitted comment, an interactive widget or gadget, or any other independent item of content.

As MDN states, an `<article>` is a _composition_ of presumably _different_ sibling elements under the `<article>` parent. Since the web is _mostly_ a vertical medium, this probably means that we are using `display: block` (ruling out the use of `gap`) and desperately looking for a means to space out our elements.

Enter: CSS margins, the direct child selector, and (an absolute classic) string manipulation.


## Using margins responsibly

I think that margins should be used in a way that they are in some way a property of the collection of elements. I’ll cover one way to do this when you are using tailwind and components, and one way for when you are using CSS selectors. 


### Tailwind

```html
<article>
    <MyHeadingComponent extraClasses=”mb-1” />
    <MyParagraphComponent extraClasses=”mb-2” />
    <MyParagraphComponent extraClasses=”mb-2” />
</article>
```

A solution that looks a bit like the above seems to fit the bill. The margins on the individual components remain visible and obvious in the parent component.

There are many solutions to the problem of getting additional classes into an already defined component using tailwind. Since I am saying that defining margins on individual elements is a poor choice, simply defining a prop for incoming classes and concatenating them onto the end of the root element’s classes shouldn’t generate any conflicts. Of course, there is also the amazing <code>[tailwind-merge](https://github.com/dcastil/tailwind-merge)</code> that will resolve any class conflicts if there happen to be any.


### Vanilla CSS

```css
article > h2 {
    margin-bottom: 0.5rem;
}

article > p {
    margin-bottom: 1rem;
}
```

Using the [direct child selector](https://developer.mozilla.org/en-US/docs/Web/CSS/Child_combinator) we can have margins defined relative to the parent element. This achieves a similar result, we can have our margins for a composition of elements defined in one place. This might look a bit nicer using [SASS](https://sass-lang.com/).

```sass
article {
    & > h2 {
        margin-bottom: 0.5rem;
    }
    & > p {
        margin-bottom: 1rem;
    }
}

```

Although, when using technology like styled components, I wish I could do this:

```css
article > MyCustomComponent {
    margin-bottom: 0.5rem;
}
```

## Conclusion

CSS margins are a source of much debate and frustration for developers. Banishing margins altogether from our projects also doesn’t seem like a particularly functional solution. Instead, we might try our best to define margins minimally: relative to the parent of a collection of different child elements. 

