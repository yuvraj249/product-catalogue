/* eslint-disable no-undef */

describe('Categories Page - Supplier Admin Restrictions', () => {

  before(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
    cy.seedCategories()
    cy.seedSuppliers()
    cy.seedUsers()
  })

  beforeEach(() => {
    cy.loginAsSupplierAdmin()
    cy.visit('/admin/categories')
  })

  it('can view categories list', () => {
    cy.get('[data-cy="categories-title"]').should('contain', 'Categories')
    cy.get('[data-cy="data-table"]').should('be.visible')
  })

  it('does not see Add Category button', () => {
    cy.get('[data-cy="add-category-btn"]').should('not.exist')
  })

  it('does not see Edit buttons', () => {
    cy.get('[data-cy="edit-category-btn"]').should('not.exist')
  })

  it('does not see Delete buttons', () => {
    cy.get('[data-cy="delete-category-btn"]').should('not.exist')
  })

})
