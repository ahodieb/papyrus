![Papyrus](papyrus.png)

# Papyrus

Note taking cli tools that i use daily to keep track of my work

## Convention & format

* Notes are written in markdown

### Entry

**Example Day entry:**

```markdown
### Fri 2021/02/19
#### Some title for the entry #meeting @user, @user | 11:30/13:00
* Notes form the meeting
  * Additional points
#### Working on something else #random-tag | 13:30
* This note is still in progress and does not have an ending time yet
```

**Breakdown:**

```txt
###  -> level (defaults to ###)
Fri 2021/02/19  -> Date of entry
#### -> level of a subentry in the day (defaults to day level + #)
Some title for the entry -> entry title and description
#meeting -> hashtags to tag the entry, this helps with creating breakdowns at the end of the week.
@user -> convention to include names, usernames of people relevant to this entry
| 11:30/13:00 -> start/end timestamps of the entry
| 13:30 -> start timestamp of the entry in progress (not completed yet)
```

## TODO

* [ ] cli can open editor to last entry position
* [ ] cli generates weekly status
* [ ] cli generates stats
* [ ] add a default command that does the correct thing by default
  * [ ] add new entry in current day
  * [ ] close off the past entry timestamp
* [ ] change defined structure to previous section (instead of a header per entry)
* [ ] Implement mechanism to reformat my old journals into new format
  * Import notes from one-note
  * Import notes from quip
* [ ] port my shell scripts to python
* [ ] add installation instructions
* [ ] Generate weekly report of what was completed
* [ ] cli to add a new entry to the todo section
* [ ] rotate journal monthly ?
* [ ] reminder to add entry to the journal 
* [ ] cli runs a server to display stats, and other reports (UI)
* [ ] cli auto closes time frames, based on some rules, or some interactive wizard/ ui + server
* [ ] Smarter rotation by completion of goals? 
* [ ] add build/development instructions ( i don't expect external contributions, but could be useful for me updating the code in 6 months)
* [ ] setup auto build and test actions on github (... just want to try it)
* [ ] add build badges (.. looks good :D)
* [ ] Refactor the format file into a module
* [ ] Setup pyEnv or virtual (requirements.txt) or similar
* [ ] round time to 5 minute blocks ?

## (Bad/Crazy/Maybe?) Ideas

* [ ] Change the structure to be a tree like, entries have sub entries, ...  (sounds like too much generalization for no much value here)
* [ ] all sections and entries have timestamps  (no good reason for this now, the only format used assumes a hierarchy of date and time entries)
* [ ] Could we present smarter way of modifying files
* [ ] in memory files with papyrus handling every section separately to a septate file ? but allows notes to read all sections in the same way (sounds too complicated for what i need, and does not add much value)
