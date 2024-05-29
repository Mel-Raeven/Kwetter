describe('End-to-End Test', () => {
  it('Logs in, navigates to dashboard, posts a message', () => {
    // Visit the homepage
    cy.visit('http://localhost:5173/')

    // Verify that the homepage contains a login button
    cy.contains('Login').should('be.visible').click()

    // Verify that the user is redirected to the Cognito Hosted UI
    cy.url().should('include', 'cognito-hosted-domain') // Replace with the URL of your Cognito Hosted UI

    // Fill in username and password and click on the login button
    cy.get('input[name="username"]').type('exampleUser') // Replace with your test username
    cy.get('input[name="password"]').type('examplePassword') // Replace with your test password
    cy.get('button[type="submit"]').click()

    // Verify that the user is redirected to the dashboard
    cy.url().should('include', '/dashboard')

    // Verify that the dashboard contains a textarea
    cy.get('textarea[name="message"]').should('be.visible').type('Test message') // Type a test message

    // Click on the post message button
    cy.get('button[type="submit"]').click()

    // Verify that the message is posted successfully
    cy.contains('Message posted successfully').should('be.visible')
  })
})

