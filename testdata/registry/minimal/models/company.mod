name: Company
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
    attributes:
      - format=company
  TaxID:
    type: String
    attributes:
      - regex=^\d{2}-\d{7}$
identifiers:
  primary: ID
  name: Name