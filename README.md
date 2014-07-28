# pongo2-addons

[![Build Status](https://travis-ci.org/flosch/pongo2-addons.svg?branch=master)](https://travis-ci.org/flosch/pongo2-addons)
[![GitTip](http://img.shields.io/badge/gittip-support%20pongo-brightgreen.svg)](https://www.gittip.com/flosch/)

Official filter and tag add-ons for [pongo2](https://github.com/flosch/pongo2). Since this package uses 
3rd-party-libraries, it's in its own package.

## How to install and use

Install via `go get -u github.com/flosch/pongo2-addons`. All dependencies will be automatically fetched and installed.

Simply add the following import line **after** importing pongo2:

    _ "github.com/flosch/pongo2-addons"

All additional filters/tags will be registered automatically.

## Addons

### Filters

  - Regulars
     - **[filesizeformat](https://docs.djangoproject.com/en/dev/ref/templates/builtins/#filesizeformat)** (human-readable filesize; takes bytes as input)
     - **[slugify](https://docs.djangoproject.com/en/dev/ref/templates/builtins/#slugify)** (creates a slug for a given input)

  - Markup
     - **markdown** (parses markdown text and outputs HTML; **hint**: use the **safe**-filter to make the output not being escaped)
  
  - Humanize
     - **[intcomma](https://docs.djangoproject.com/en/dev/ref/contrib/humanize/#intcomma)** (put decimal marks into the number)
     - **[ordinal](https://docs.djangoproject.com/en/dev/ref/contrib/humanize/#ordinal)** (convert integer to its ordinal as string)
     - **[timesince](https://docs.djangoproject.com/en/dev/ref/templates/builtins/#timesince)/[timeuntil](https://docs.djangoproject.com/en/1.6/ref/templates/builtins/#timeuntil)/[naturaltime](https://docs.djangoproject.com/en/dev/ref/contrib/humanize/#naturaltime)** (human-readable time [duration] indicator)

### Tags

(nothing yet)

## TODO

 - Support i18n/i10n