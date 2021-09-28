#!/bin/bash
go get -u github.com/liberty239/cassiopaea-tools/./...
cass-src fetch-all .build/repo
cass-gen epub .build/repo .build/sessions.epub
cass-gen html .build/repo .build/sessions.html
