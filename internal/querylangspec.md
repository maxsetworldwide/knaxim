# Knaxim Search Query Structure

A query has two parts a context and a matching condition represented in json.

```java
{
  "context": <context_value>,  
  "match": <match_value>
}
```

## Context

The search context represents a set of files. The context value in a query defines that initial set that the match condition filters over. The initial data type of the context value determines how it should be interpreted.

- Array  
An array context value means that each element of the array should be interpreted as a context value and the resulting file set is the combination of all the files from each element

- Object  
The context object has two required fields: "type" and "id". "type" determines how the object is interpreted and what additional fields the object may have. "id" is the primary identifier for that type.

  - "owner"  
  the id is the owner id value as a string. By default it searches both owned and viewable files by the owner. there is an optional field of "only" that has 2 valid values of "owned" and "viewable" which limits to owned or viewable files respectively.

  - "file"  
  The id is the id of the file. This context represents a single file

- String  
A string is short hand for the id value of the object with type "owner".
`"aaaaa"` becomes  
```json
{  
  "type": "owner",
  "id": "aaaaa"
}
```

## Match

The match value represents the filter condition for the context. The data type of the match field if the first determiner of how it is interpreted.

- Array  
Each element of an array will be interpreted as a match value and collectively only files that match every element of the array will match the whole value

- Object  
  - Required
    - "tagtype"  
    can be a string which is interpreted to be one of the tag type values, a number which is cast to a tag type, or an array of tagtype values which overall is the combination of all tag types in the array.
    - "word"  
    what to match within the type of tag. this matching always ignores case.
  - Optional
    - "regex"  
    indicates that the "word" is the be interpreted as a regular expression to match tags against instead of just equality. As long as regex is present and not null the "word" will be interpreted as a regular expression.

- String  
will be interpreted as a regular expression seartching the content "tagtype".  
`cat` becomes  
```json
{
  "tagtype": "content",
  "word": "cat",
  "regex": true
}
```

The current types of tags are:
- content
- topic
- action
- process
- resource
- user
- date
- name
