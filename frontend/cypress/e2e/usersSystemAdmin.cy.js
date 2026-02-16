/* eslint-disable no-undef */

describe('Users Page - System Admin', () => {

  const userName = "Test User"
  const updatedName = "Updated User"
  const userEmail = "testuser@test.com"
  const updatedEmail = "updateduser@test.com"
  const userPassword = "Test@1234"

  before(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
    cy.seedSuppliers()
  })

  beforeEach(() => {
    cy.loginAsSystemAdmin()
    cy.visit('/admin/users')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads Users page', () => {
    cy.contains('Users').should('exist')
    cy.get('[data-cy="add-user-btn"]').should('exist')
  })

  it('opens Add User modal', () => {
    cy.get('[data-cy="add-user-btn"]').click()
    cy.get('[data-cy="user-modal-title"]').should('contain.text', 'Add User')
  })

  it('blocks empty form submission', () => {
    cy.get('[data-cy="add-user-btn"]').click()
    cy.get('[data-cy="user-submit-btn"]').click()

    cy.contains('User name is required').should('exist')
  })

  it('blocks invalid email format', () => {
    cy.get('[data-cy="add-user-btn"]').click()

    cy.get('[data-cy="user-name"]').type('John')
    cy.get('[data-cy="user-email"]').type('invalidemail')
    cy.get('[data-cy="user-password"]').type(userPassword)

    cy.get('[data-cy="user-submit-btn"]').click()
    cy.contains('Invalid email format (e.g., user@example.com)').should('exist')
  })

  it('blocks weak password', () => {
    cy.get('[data-cy="add-user-btn"]').click()

    cy.get('[data-cy="user-name"]').type('John')
    cy.get('[data-cy="user-email"]').type('john@test.com')
    cy.get('[data-cy="user-password"]').type('weak')

    cy.get('[data-cy="user-submit-btn"]').click()
    cy.contains('Password must be at least 8 characters long').should('exist')
  })

  it('requires supplier selection', () => {
    cy.get('[data-cy="add-user-btn"]').click()

    cy.get('[data-cy="user-name"]').type('John')
    cy.get('[data-cy="user-email"]').type('john@test.com')
    cy.get('[data-cy="user-password"]').type(userPassword)

    cy.get('[data-cy="user-submit-btn"]').click()
    cy.contains('Please select a valid supplier').should('exist')
  })

  it('creates a new supplier admin user', () => {
  cy.get('[data-cy="add-user-btn"]').click()

  cy.get('[data-cy="user-name"]').type(userName)
  cy.get('[data-cy="user-email"]').type(userEmail)

  cy.get('.react-select__control').click({ force: true })
  cy.get('.react-select__menu')
    .contains('Supplier')
    .first()
    .click()

  cy.get('[data-cy="user-password"]').type(userPassword)

  cy.get('[data-cy="user-submit-btn"]').click()

  cy.contains(userName, { timeout: 8000 }).should('exist')
})


  it('edits an existing user', () => {
  cy.contains('td', userName)
    .parents('tr')
    .find('[data-cy="edit-user-btn"]')
    .click()

  cy.get('[data-cy="user-name"]').clear().type(updatedName)
  cy.get('[data-cy="user-email"]').clear().type(updatedEmail)

  cy.get('[data-cy="user-submit-btn"]').click()

  cy.contains(updatedName).should('exist')
})

  it('closes modal using Cancel', () => {
    cy.get('[data-cy="add-user-btn"]').click()
    cy.get('[data-cy="user-cancel-btn"]').click()

    cy.get('[data-cy="user-modal-title"]').should('not.exist')
  })

  it('deletes the user', () => {
    cy.contains('td', updatedName)
      .parents('tr')
      .find('[data-cy="delete-user-btn"]')
      .click()

    cy.contains(updatedName).should('not.exist')
  })

  it('searches users', () => {
    cy.get('[data-cy="user-search"]').type('Supplier')
    cy.contains('Supplier').should('exist')
  })

})
