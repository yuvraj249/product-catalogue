/* eslint-disable no-undef */

describe('Dashboard - Supplier Admin', () => {

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
    cy.loginAsSupplierAdmin()
    cy.visit('/admin/dashboard')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads dashboard', () => {
    cy.get('[data-cy="dashboard-title"]').should('contain', 'Dashboard')
  })

  it('shows all stat cards', () => {
    cy.get('[data-cy="dashboard-card-total-products"]').should('exist')
    cy.get('[data-cy="dashboard-card-suppliers"]').should('exist')
    cy.get('[data-cy="dashboard-card-supplier-admins"]').should('exist')
    cy.get('[data-cy="dashboard-card-categories"]').should('exist')

  })

  it('shows only company low stock products', () => {
    cy.get('[data-cy="low-stock-table"]').contains('Laptop')
  })

  it('cannot access users via dashboard card', () => {
  cy.get('[data-cy="dashboard-card-supplier-admins"]').click()
  cy.contains('access denied').should('exist')
})

  it('navigates to products', () => {
    cy.get('[data-cy="dashboard-card-total-products"]').click()
    cy.url().should('include', '/admin/products')
  })

})
