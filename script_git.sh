find . -name '*.DS_Store' -type f -delete
git add *
git commit -m "update"
git push -u origin master
