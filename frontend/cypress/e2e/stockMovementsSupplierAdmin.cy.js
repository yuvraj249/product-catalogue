/* eslint-disable no-undef */

describe('Stock Movements - Supplier Admin', () => {

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
    cy.visit('/admin/stock-movements')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads Stock page', () => {
    cy.get('[data-cy="stock-title"]').should('contain', 'Stock')
    cy.get('[data-cy="data-table"]').should('be.visible')
    cy.get('[data-cy="add-stock-btn"]').should('be.visible')
  })

  it('opens Add Stock modal', () => {
  cy.get('[data-cy="add-stock-btn"]').click()
  cy.get('[data-cy="stock-modal-title"]').should('be.visible')
 })


  it('blocks empty form submission', () => {
    cy.get('[data-cy="add-stock-btn"]').click()
    cy.get('[data-cy="stock-submit-btn"]').click()
    cy.contains('Please select a valid product').should('exist')
  })

  it('blocks invalid quantity (negative)', () => {
  cy.get('[data-cy="add-stock-btn"]').click()

  cy.contains('Product *')
    .parent()
    .find('.react-select__control')
    .click()

  cy.get('.react-select__menu')
    .contains('Laptop')
    .click()

  cy.get('[data-cy="stock-quantity"]').type('-5')
  cy.get('[data-cy="stock-submit-btn"]').click()

  cy.contains('positive number').should('exist')
})


  it('blocks decimal quantity', () => {
  cy.get('[data-cy="add-stock-btn"]').click()

  cy.contains('Product *')
    .parent()
    .find('.react-select__control')
    .click()

  cy.get('.react-select__menu')
    .contains('Laptop')
    .click()

  cy.get('[data-cy="stock-quantity"]').type('2.5')
  cy.get('[data-cy="stock-submit-btn"]').click()

  cy.contains('Quantity must be a whole number').should('exist')
})


  it('blocks extremely large quantity', () => {
  cy.get('[data-cy="add-stock-btn"]').click()

  cy.contains('Product *')
    .parent()
    .find('.react-select__control')
    .click()

  cy.get('.react-select__menu')
    .contains('Laptop')
    .click()

  cy.get('[data-cy="stock-quantity"]').type('9999999')
  cy.get('[data-cy="stock-submit-btn"]').click()

  cy.contains('Quantity too large').should('exist')
})


  it('creates a valid stock movement', () => {
  cy.get('[data-cy="add-stock-btn"]').click()

  cy.contains('Product *')
    .parent()
    .find('.react-select__control')
    .click()

  cy.get('.react-select__menu')
    .contains('Laptop')
    .click()

  cy.get('[data-cy="stock-quantity"]').type('5')

  cy.contains('Movement Type *')
    .parent()
    .find('.react-select__control')
    .click()

  cy.get('.react-select__menu')
    .contains('IN')
    .click()

  cy.get('[data-cy="stock-reason"]').type('New delivery')
  cy.get('[data-cy="stock-submit-btn"]').click()
  cy.contains('New delivery', { timeout: 8000 }).should('exist')
})


  it('edits a stock movement', () => {
  cy.get('[data-cy="edit-stock-btn"]').first().click()

  cy.get('[data-cy="stock-quantity"]').clear().type('9')
  cy.get('[data-cy="stock-reason"]').clear().type('Updated stock')

  cy.get('[data-cy="stock-submit-btn"]').click()

  cy.contains('Updated stock', { timeout: 8000 }).should('exist')
})


  it('searches stock movements', () => {
    cy.get('[data-cy="stock-search"]').type('Updated')
    cy.contains('Updated stock').should('exist')
  })
  

  it('deletes a stock movement', () => {
    cy.get('[data-cy="delete-stock-btn"]').first().click()
    cy.contains('Updated stock').should('not.exist')
  })

})
