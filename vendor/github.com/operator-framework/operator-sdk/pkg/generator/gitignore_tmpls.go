// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

const projectGitignoreTmpl = `
# Temporary Build Files
tmp/_output
tmp/_test


# Created by https://www.gitignore.io/api/go,vim,emacs,visualstudiocode

### Emacs ###
# -*- mode: gitignore; -*-
*~
\#*\#
/.emacs.desktop
/.emacs.desktop.lock
*.elc
auto-save-list
tramp
.\#*

# Org-mode
.org-id-locations
*_archive

# flymake-mode
*_flymake.*

# eshell files
/eshell/history
/eshell/lastdir

# elpa packages
/elpa/

# reftex files
*.rel

# AUCTeX auto folder
/auto/

# cask packages
.cask/
dist/

# Flycheck
flycheck_*.el

# server auth directory
/server/

# projectiles files
.projectile
projectile-bookmarks.eld

# directory configuration
.dir-locals.el

# saveplace
places

# url cache
url/cache/

# cedet
ede-projects.el

# smex
smex-items

# company-statistics
company-statistics-cache.el

# anaconda-mode
anaconda-mode/

### Go ###
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, build with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

### Vim ###
# swap
.sw[a-p]
.*.sw[a-p]
# session
Session.vim
# temporary
.netrwhist
# auto-generated tag files
tags

### VisualStudioCode ###
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
.history


# End of https://www.gitignore.io/api/go,vim,emacs,visualstudiocode
`
