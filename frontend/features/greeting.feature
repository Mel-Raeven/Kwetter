Feature: Post Message

Scenario: User posts a message
  Given the user is logged in with a valid ID token
  When the user posts "Hello, world!"
  Then the message should be sent

