---
title: Making Spotify Wrapped for iTunes.xml files using TreeSitter ... in 2024
date: 2024-02-22
slug: itunes-xml
description: >
  iTunes for Windows and MacOS offers data export in the form of XML files. As part of my major computer science assessment, I build a parser for these files to extract listening statistics.
---
Library export was a feature added to iTunes that allowed developers to build integrations with it. Now that the iPod has had its resurgence, many find themselves in the position of using this two decade old digital music library. 

Given that streaming services like Apple Music and Spotify are incentivised by their monthly paying customers to provide an engaging and innovative experience, iTunes' offerings are lacking in surprisingly few ways. The app sports smart-playlists and even 'Genius', a recommender system providing smart shuffling, mixes, and playlists similar to the AI features of today.

However, there has been one major innovation in the digital music library space that iTunes simply does not feature. The annual statistical summary. I knew I had to take iTunes' ancient XML export feature and turn it into a Spotify Wrapped inspired experience for all of the iPod users of today.

An application such as this one is expected by users to run in their web browser. So what language is better used to parse an exported iTunes.XML file than the web's flagship language: JavaScript.

I considered the DOMParser API and XML parsers more broadly, but concluded that building my own would be more fun. These options would also only bring me slightly closer to a JavaScript object of a user's library, since they do not understand how an iTunes.XML file specifically is formatted.

Parsers are a well researched area of computer science. The first step in most parsing algorithms is to extract tokens (things like the word 'if' in many programming languages or '<' in an XML file). 'Lexical analysis' is performed on an incoming stream of characters to notice these tokens as they come in. 

I built a first attempt [lexical analyser in JavaScript](https://github.com/will-lol/ComputerScienceIA/blob/42c8c9ab3665056f6382edb825b6654ed52c9e70/src/lib/util/TokenReader.ts) and it had worked on some small test files I gave it. That was until I gave it an entire 1M iTunes.xml file and the thing ran out of memory. Perhaps I had proven [Atwood's law](https://en.wikipedia.org/wiki/Jeff_Atwood).

Me and my horrible code definitely bear some responsibility for this terrible performance, but I felt like writing this parser in JavaScript was simply a poor decision.

WASM is a web technology attracting a persistent crowd of developers (including myself) claiming it will replace JavaScript in <span class="font-serif">𝑥</span> amount of years. Used in [Atom](https://atom-editor.cc/) (RIP), TreeSitter is a parser generator library that boasts speed in both parsing and the writing of parsers. It also has an option to export WASM binaries. 

And so, I wrote my own TreeSitter grammar for the iTunes.XML file.

```
module.exports = grammar({
  name: 'iTunesXML', 
  rules: {
    source_file: $ => repeat(choice($.doctype, $.xml_declaration, $.plist)),
    doctype: $ => seq('<!', /[Dd][Oo][Cc][Tt][Yy][Pp][Ee]/, /[^>]+/, '>'),
    xml_declaration: $ => seq('<?', /[Xx][Mm][Ll]/, /[^?]+/, '?>'),
    plist: $ => seq($._plistStart, repeat($._expression), $._plistEnd),
    _plistStart: $ => seq('<', /[Pp][Ll][Ii][Ss][Tt]/, /[^>]+/, '>'),
    _plistEnd: $ => seq('</', /[Pp][Ll][Ii][Ss][Tt]/, '>'),
    obj: $ => seq('<dict>', repeat($.item), '</dict>'),
    item: $ => seq($.key, $._expression),
    key: $ => seq('<key>', $.text, '</key>'),
    text: $ => /[^<>s]([^<>]*[^<>s])?/,
    _expression: $ => choice($.obj, $.array, $.integer, $.real, $._boolean, $.data, $.date, $.string),
    array: $ => seq('<array>', repeat($._expression), '</array>'),
    integer: $ => seq('<integer>', $.int, '</integer>'),
    int: $ => /[0-9]+/,
    string: $ => seq('<string>', optional($.text), '</string>'),
    real: $ => seq('<real>', $.float, '</real>'),
    float: $ => /(-|+)?d+.?d*(E(-|+)[0-9]+)?/,
    _boolean: $ => choice($.true, $.false),
    true: $ => '<true/>',
    false: $ => '<false/>',
    data: $ => seq('<data>', $.base64, '</data>'),
    date: $ => seq('<date>', $.iso8601, '</date>'), 
    iso8601: $ => /d{4}(-dd(-dd(Tdd:dd(:dd)?(.d+)?(([+-]dd:dd)|Z)?)?)?)?/,
    base64: $ => /(?:[A-Za-z0-9s+/]{4})*(?:[A-Za-z0-9s+/]{2}==|[A-Za-z0-9s+/]{3}=)?/
  }
})
```

My grammar certainly isn't perfect. This is my first time writing one!

Each 'rule' in 'rules' defines a construct that exists in the file. Rules prefixed with '\_' will not appear in the final syntax tree. You usually use these for wrappers like `_expression` or in `_boolean`, where I didn't want a `_boolean` construct surrounding the actual value of `true` or `false`. An [example iTunes.XML](https://github.com/will-lol/ComputerScienceIA/blob/main/tests/test.xml) is available.

Back in JavaScript, I now needed to process the result given by TreeSitter into a useful JavaScript object. 

```ts
// This code maps a TreeSitter array of all of the tracks in our XML into an array of JavaScript song objects.
const songs = tracksDict
  .descendantsOfType('obj')
  .slice(1)
  .map((item): song => {
    let name: string | undefined,
        artist: string | undefined,
        album: string | undefined,
        genre: string | undefined;
    let time: number | undefined,
        playCount: number | undefined,
        skipCount: number | undefined,
        rating: number | undefined;

    // We are going to walk down the syntax tree, looking for any of the variables relevant to the JavaScript song object.
    const cursor = item.walk();
    cursor.gotoFirstChild(); // The current node is a 'dict', we would like to enter the dict and start looping over its siblings.
    while (cursor.gotoNextSibling()) {
      if (cursor.nodeType == 'item') {
        const key = cursor.currentNode().namedChild(0); // The key is the first named child of each 'item' node in a 'dict'.
        if (key == null) {
          throw 'no key in item';
        }
        const keyName = getKeyName(key);

        const dataNode = key.nextSibling; // The 'data' node is the following sibling of the key
        if (dataNode == null) {
          throw 'no data node in item';
        }
        // We just switch on the key to see if it is relevant. If so, we just get the value of its dataNode and set it to its associated variable declared earlier.
        switch (keyName) {
          case 'Name':
            name = getAndParseKeyString(dataNode);
            break;
          case 'Artist':
            artist = getAndParseKeyString(dataNode);
            break;
          case 'Album':
            album = getAndParseKeyString(dataNode);
            break;
          case 'Genre':
            genre = getAndParseKeyString(dataNode);
            break;
          case 'Total Time':
            time = getAndParseKeyNumber(dataNode);
            break;
          case 'Play Count':
            playCount = getAndParseKeyNumber(dataNode);
            break;
          case 'Skip Count':
            skipCount = getAndParseKeyNumber(dataNode);
            break;
          case 'Rating':
            rating = getAndParseKeyNumber(dataNode);
            break;
        }
      }
    }
    // Now we can assemble our song object.
    return {
      name: name,
      artist: artist,
      album: album,
      genre: genre,
      time: time,
      playCount: playCount,
      skipCount: skipCount,
      rating: rating
    };
});
```
The [full source code](https://github.com/will-lol/ComputerScienceIA/blob/main/src/routes/parseWorker.ts) is available. If you are interested to try my iTunes.XML parser and generate some statistics for yourself, it is [hosted on Vercel](https://computer-science-ia.vercel.app/).
