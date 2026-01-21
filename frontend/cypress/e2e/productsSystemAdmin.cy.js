/* eslint-disable no-undef */

describe('Products Page - System Admin', () => {

  before(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
    cy.seedCategories()
    cy.seedSuppliers()
    cy.seedUsers()
    cy.loginAsSupplierAdmin()
    cy.seedProducts()
  })

  beforeEach(() => {
    cy.loginAsSystemAdmin()
    cy.visit('/admin/products')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads Products page', () => {
    cy.get('[data-cy="products-title"]').should('contain', 'Products')
    cy.get('[data-cy="products-table"]').should('be.visible')
  })

  it('does NOT show Add Product button', () => {
    cy.get('[data-cy="add-product-btn"]').should('not.exist')
  })

  it('does NOT show Edit buttons', () => {
    cy.get('[data-cy="edit-product-btn"]').should('not.exist')
  })

  it('does NOT show Delete buttons', () => {
    cy.get('[data-cy="delete-product-btn"]').should('not.exist')
  })

  it('cannot open product modal via UI', () => {
    cy.get('[data-cy="product-modal-title"]').should('not.exist')
  })

  it('can search products', () => {
    cy.get('[data-cy="product-search"]').type('Laptop')

    cy.get('[data-cy="products-table"]', { timeout: 8000 })
      .should('contain.text', 'Laptop')
  })

})
