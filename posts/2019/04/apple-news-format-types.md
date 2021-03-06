title:Apple News Format Types
date:2019-04-29
excerpt:A library of expressions for working with the Apple News Format in publishing
image:anflib.png
---
<small>2019-04-29</small>

<div class="blog__post-header-image" style="background-image: url(/static/images/anflib.png)">&nbsp;
</div>

<a class="github-button" href="https://github.com/robertfairley/apple-news-format" data-icon="octicon-star" data-size="large" aria-label="Star robertfairley/apple-news-format on GitHub">Star</a> <a class="github-button" href="https://github.com/robertfairley/apple-news-format/fork" data-icon="octicon-repo-forked" data-size="large" aria-label="Fork robertfairley/apple-news-format on GitHub">Fork</a>
<hr/>

**Recently I've been involved in a project** converting print design files to the newly [available in Canada] Apple News Format.

My role in that endeavour has been in designing and development of a service that composes a proprietary input format into standard Apple News Format articles fit for publication to Apple's "plus" service that recently replaced [Texture](https://techcrunch.com/2019/03/29/apple-to-close-texture-on-may-28-following-launch-of-apple-news/).

In the progress of that project I encountered a lack of community support for publishing to or developing for the platform. There were a few groups who implemented some access to the <abbr class="tooltip" data-tooltip="Apple News API">API</abbr> in a couple of languages, but fewer that kept up with the recent changes and the recent [announcement](https://www.apple.com/ca/apple-events/march-2019/).

Because of my necessary work with the current <abbr class="tooltip"  data-tooltip="Apple News Format">ANF</abbr> <abbr class="tooltip" data-tooltip="Apple News API">API</abbr> and the potential for some needed agreement between disparate application interfaces in those proprietary systems I decided to implement all of the <abbr class="tooltip" data-tooltip="Apple News Format">ANF</abbr> TypeScript types and publish them as an installable package in case anyone else might find them helpful.

For now you can find that available by cloning the repository or installing the <abbr class="tooltip" data-tooltip="Node Package Manager">NPM</abbr> package.

<hr/>

*Clone Repository*

```sh
git clone https://github.com/robertfairley/apple-news-format.git
```

*Import Package*

```sh
npm i apple-news-format
```
<hr/>

I definitely welcome any help anybody might want to provide, if you find the code useful.

Eventually I might decide to add implementations for each of the provided interfaces, but for now I'm leaving implementation up to the package consumer. That way the <abbr class="tooltip" data-tooltip="Apple News Format">ANF</abbr> contracts are provided, but all opinionated interference ends there.
