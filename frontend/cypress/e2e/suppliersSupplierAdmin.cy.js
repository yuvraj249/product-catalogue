/* eslint-disable no-undef */

describe('Suppliers Page - Supplier Admin Restrictions', () => {

  before(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
    cy.seedSuppliers()
    cy.seedUsers()
  })


  beforeEach(() => {
    cy.loginAsSupplierAdmin()
    cy.visit('/admin/suppliers')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('can view suppliers list', () => {
    cy.contains('Suppliers').should('be.visible')
    cy.get('[data-cy="suppliers-table"]').should('be.visible')
  })

  it('does not see Add Supplier button', () => {
    cy.get('[data-cy="add-supplier-btn"]').should('not.exist')
  })

  it('does not see Edit buttons', () => {
    cy.get('[data-cy="edit-supplier-btn"]').should('not.exist')
  })

  it('does not see Delete buttons', () => {
    cy.get('[data-cy="delete-supplier-btn"]').should('not.exist')
  })

})
