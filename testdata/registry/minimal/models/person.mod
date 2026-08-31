name: Person
fields:
  ID:
    type: AutoIncrement
  FirstName:
    type: String
    attributes:
      - format=firstName
  LastName:
    type: String
    attributes:
      - format=lastName
  Email:
    type: String
    attributes:
      - format=email
  Nationality:
    type: Nationality
  Bio:
    type: String
    attributes:
      - optional
      - maxLength=200
identifiers:
  primary: ID
related:
  Company:
    type: ForOne