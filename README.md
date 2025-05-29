# Knowledge bank management system

This is my accepted proposal for the knowledge bank management system.
The system is designed to help users manage their knowledge base effectively,
allowing them to create, update, and delete entries,
as well as search for specific information.

> This database architecture was accepted at Impel

This structure of having each separate tables allows for easy management and scalability.

Theoretically there wouldn't be duplicated data in "helper" tables.
Only place of coupling is the `entries` table.

## Entries Table

```text
-------------------------------------------------------------------------
| Id | DealerId | GroupId | ProductId | TopicId | QuestionId | AnswerId |
-------------------------------------------------------------------------
```

Groups and Products are predefined so basically an enum
There are some Topics and Questions that are "defaults"
Meaning, whenever new dealer is created these defaults are assigned to the dealer.

For example:

```text
---------------------------------------------------------------------------------
| Id | DealerId | GroupId | ProductId |     TopicId |     QuestionId | AnswerId |
---------------------------------------------------------------------------------
|  1 | dealer-1 |       1 |         1 | def-topic-1 | def-question-1 |     null |
---------------------------------------------------------------------------------
|  2 | dealer-1 |       1 |         1 | def-topic-1 | def-question-2 |     null |
---------------------------------------------------------------------------------
```

> Based on these two entries we can tell that our dealer has one topic with two questions

### Scenarios

Now i will go over the cases and scenarios of knowledge bank usage

#### When question is answered

1. new answer is created in `Answers` table
2. `AnswerId` is updated in the `Entries` table

#### When new question is created

- it's required to create new question without an answer

1. new question is created in `Questions` table
2. new answer is created in `Answers` table
3. new entry is created in `Entries` table with
    - `QuestionId` from the step 1.
    - `AnswerId` from the step 2.
    - rest values are passed from the front-end

#### When new topic is created

> NOTE: if the topic isn't used instantly it will disappear from the front-end
> and garbage collector will remove it after some time

- these topics have special property `Custom` set to `true`
- with this property user can delete the topic later

1. first we check whether we have this Topic in `Topics` table:
    - if we do, we return it's id
    - if not, we add new row to `Topics` table and return it's id

#### Copying entries

- We get the source path and destination path
- we copy all entries from source to destination
