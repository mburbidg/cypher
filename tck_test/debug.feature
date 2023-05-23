Feature: Debug - Debug test cases

  Scenario: [18] Fail when creating a relationship without a type
    Given any graph
    When executing query:
      """
      CREATE ()-->()
      """
    Then a SyntaxError should be raised at compile time: NoSingleRelationshipType
