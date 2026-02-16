/* eslint-disable no-undef */
describe('Auth Guard', () => {

  it('redirects unauthenticated user from /admin/categories', () => {
    cy.clearLocalStorage()
    cy.visit('/admin/categories')
    cy.url().should('eq', Cypress.config().baseUrl + '/')
  })

})

