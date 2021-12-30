# kape

To test out, run:

```
go run main.go puc19.gbk
```

### TODO
This is not usable currently. To implement:
1. Ability to select sequence (vim-like)
2. Proper color high-lighting (and sane defaults)
3. Good information viewing for feature list (maybe change to grid)

After these, the plan is:
1. Implement ability to change sequence
2. Feature addition and import
3. io.Writer and stream implementation (to scroll through genomes)
4. Find capability

### Hotkeys
```
escape: get out of particular view (sequence highlight, etc)
a: switch to action view
s: switch to sequence view
f: switch to feature view
g: switch to genbank view with default editor (vim)

h: left
j: down
k: up
l: right

### Sequence view
t: toggle between sequence view (ape-like) and feature view (snapgene-like)

b: Jump 10 bases backwards
e: Jump 10 bases forwards

v: highlight sequence
- y: copy sequence
- c: copy (to clipboard)
- d: delete sequence

p: paste sequence
x: delete base
i: insert base

```

