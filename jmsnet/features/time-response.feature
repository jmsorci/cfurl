Feature: time website response

   Scenario: time responsive site
     When I specify a valid URL like "http://www.cnn.com"
     And I specify a timeout of 2000 ms
     And I measure response time
     Then positive fetch time less than timeout is reported