# CERES

## Methods available

### AddEntry

Request type : `add_entry`

Arguments : 

 - **word**, the word used to add this onto the database
    - Formulated as `{word}`
    - required
 - **isType**, whether the word refers to a type or a specific request
    - Formulated as `{istype=(y|n)}`
    - required
 - **parent**, what that word is a hyponym of (for example )
    - formulated as `--parent {parent_word} {parent_index}`
    - **not** required
 - **grammar group**
    - formulated as `--ggroup {ggroup_name}`
    - **not** required
