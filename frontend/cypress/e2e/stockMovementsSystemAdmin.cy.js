/* eslint-disable no-undef */

describe('Stock Movements - System Admin', () => {

  before(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()

    cy.seedCategories()
    cy.seedSuppliers()
    cy.seedUsers()

    cy.loginAsSupplierAdmin()
    cy.seedProducts()
    cy.seedStock()
  })

  beforeEach(() => {
    cy.loginAsSystemAdmin()
    cy.visit('/admin/stock-movements')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads Stock page', () => {
    cy.get('[data-cy="stock-title"]').should('contain', 'Stock')
    cy.get('[data-cy="data-table"]').should('be.visible')
  })

  it('does NOT show Add Stock button', () => {
    cy.get('[data-cy="add-stock-btn"]').should('not.exist')
  })

  it('does NOT show Edit buttons', () => {
    cy.get('[data-cy="edit-stock-btn"]').should('not.exist')
  })

  it('does NOT show Delete buttons', () => {
    cy.get('[data-cy="delete-stock-btn"]').should('not.exist')
  })

  it('can search stock movements', () => {
    cy.get('[data-cy="stock-search"]').type('Initial')
    cy.contains('Initial stock').should('exist')
  })

  it('can see User column', () => {
    cy.contains('User').should('exist')
  })

})
