/* eslint-disable no-undef */

describe('Dashboard - System Admin', () => {

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
    cy.get('[data-cy="dashboard-card-categories"]').should('exist')
    cy.get('[data-cy="dashboard-card-suppliers"]').should('exist')
    cy.get('[data-cy="dashboard-card-supplier-admins"]').should('exist')
  })

  it('navigates on card click', () => {
    cy.get('[data-cy="dashboard-card-total-products"]').click()
    cy.url().should('include', '/admin/products')
  })

  it('shows low stock section', () => {
    cy.get('[data-cy="low-stock-title"]').should('contain', 'Low Stock')
    cy.get('[data-cy="low-stock-threshold"]').should('contain', 'Threshold')
    cy.get('[data-cy="data-table"]').should('exist')
  })

  it('renders low stock products', () => {
    cy.get('[data-cy="data-table"]').contains('Laptop')
  })

  it('sorts product names', () => {
    cy.contains('Product Name').click()
  })

})
