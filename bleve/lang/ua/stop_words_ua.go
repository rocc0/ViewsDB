package ua

import (
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
)

const StopName = "stop_ua"

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/snowball/
// ` was changed to ' to allow for literal string

var UkrainianStopWords = []byte(` | From svn.tartarus.org/snowball/trunk/website/algorithms/ukrainian/stop.txt
 | This file is distributed under the BSD License.
 | See http://snowball.tartarus.org/license.php
 | Also see http://www.opensource.org/licenses/bsd-license.html
 |  - Encoding was converted to UTF-8.
 |  - This notice was added.
 |
 | NOTE: To use this file with StopFilterFactory, you must specify format="snowball"

 | a ukrainian stop word list. comments begin with vertical bar. each stop
 | word is at the start of a line.

 | this is a ranked list (commonest to rarest) of stopwords derived from
 | a large text sample.

 |

і              | and
в              | in/into
не             | not
що             | what/that
він            | he
на             | on/onto
я              | i
з              | from
із             | alternative form
як             | how
а              | milder form of 'no' (but)
то             | conjunction and form of 'that'
все            | all
вона           | she
отже           | so, thus
його           | him
але            | but
так            | yes/and
ти             | thou
до             | towards, by
біля           | around, chez
ж              | intensifier particle
ви             | you
за             | beyond, behind
би             | conditional/subj. particle
по             | up to, along
тільки         | only
її             | her
мені           | to me
було           | it was
от             | here is/are, particle
від            | away from
мене           | me
ще             | still, yet, more
ні             | no, there isnt/arent
про            | about
з              | out of
йому           | to him
тепер          | now
коли           | when
навіть         | even
ну             | so, well
раптом         | suddenly
якщо           | if
вже            | already, but homonym of 'narrower'
або            | or
ні             | neither
бути           | to be
був            | he was
його           | prepositional form of его
до             | up to
вас            | you accusative
небуть         | indef. suffix preceded by hyphen
знову          | again
вже             | already, but homonym of 'adder'
вам            | to you
сказав         | he said
адже           | particle 'after all'
там            | there
потім          | then
себе           | oneself
нічого         | nothing
їй             | to her
може           | usually with 'быть' as 'maybe'
вони           | they
тут            | here
де             | where
є              | there is/are
треба          | got to, must
нею            | prepositional form of  ей
для            | for
ми             | we
тебя           | thee
їх             | them, their
чим            | than
була           | she was
сам            | self
щоб            | in order to
без            | without
ніби           | as if
людина         | man, person, one
чого           | genitive form of 'what'
раз            | once
теж            | also
собі           | to oneself
під            | beneath
життя          | life
буде           | will be
ж              | short form of intensifer particle 'же'
тоді           | then
хто            | who
цей            | this
говов          | was saying
того           | genitive form of 'that'
тому           | for that reason
цього          | genitive form of 'this'
який           | which
зовсім         | altogether
ним            | prepositional form of 'его', 'они'
тут            | here
цьому          | prepositional form of 'этот'
один           | one
майже          | almost
мій            | my
тим            | instrumental/dative plural of 'тот', 'то'
щоб            | full form of 'in order that'
неї            | her (acc.)
здається       | it seems
зараз          | now
були           | they were
куди           | where to
навіщо         | why
сказати        | to say
всіх           | all (acc., gen. preposn. plural)
ніколи         | never
сьогодні       | today
можна          | possible, one can
при            | by
накінець       | finally
два            | two
про            | alternative form of 'о', about
інший          | another
хоч            | even
після          | after
над            | above
більше         | more
той            | that one (masc.)
через          | across, in
ці             | these
нас            | us
про            | about
всього         | in all, only, of all
них            | prepositional form of 'они' (they)
яка            | which, feminine
багато         | lots
хіба           | interrogative particle
сказала        | she said
три            | three
цю             | this, acc. fem. sing.
моя            | my, feminine
втім           | moreover, besides
добре          | good
свою           | ones own, acc. fem. sing.
цій            | oblique form of 'эта', fem. 'this'
перед          | in front of
інколи         | sometimes
краще          | better
трохи          | a little
тому           | preposn. form of 'that one'
такий          | such a one
їм             | to them
більш          | more
завжди         | always
звичайно       | of course
всю            | acc. fem. sing of 'all'
між            | between

`)

func TokenMapConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenMap, error) {
	rv := analysis.NewTokenMap()
	err := rv.LoadBytes(UkrainianStopWords)
	return rv, err
}

func init() {
	registry.RegisterTokenMap(StopName, TokenMapConstructor)
}