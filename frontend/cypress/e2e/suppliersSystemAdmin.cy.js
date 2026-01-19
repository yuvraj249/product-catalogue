/* eslint-disable no-undef */

describe('Suppliers Page - System Admin', () => {

  const supplierName = "Test Supplier"
  const updatedName = "Updated Supplier"

  beforeEach(() => {
    cy.loginAsSystemAdmin()
    cy.visit('/admin/suppliers')
  })

  it('loads suppliers page', () => {
    cy.contains('Suppliers').should('be.visible')
    cy.get('[data-cy="add-supplier-btn"]').should('be.visible')
    cy.get('[data-cy="suppliers-table"]').should('be.visible')
  })

  it('opens Add Supplier modal', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-modal-title"]').should('contain.text', 'Add Supplier')
  })


  it('blocks empty form submission', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Supplier name required')
  })

  it('blocks short supplier name', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('A')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Supplier name too short')
  })

  it('blocks long supplier name', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('A'.repeat(60))
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Supplier name too long')
  })

  it('blocks invalid characters in name', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('@@@###')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Supplier name should only contain alphabets and spaces')
  })

  it('blocks empty contact', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('Valid Name')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Contact number required')
  })

  it('blocks invalid contact format', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('Valid Name')
    cy.get('[data-cy="supplier-contact"]').type('abcd')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Contact should contain only numbers, +, -, () or spaces')
  })

  it('blocks invalid email', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('Valid Name')
    cy.get('[data-cy="supplier-contact"]').type('+919999999999')
    cy.get('[data-cy="supplier-email"]').type('invalidemail')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Invalid email format')
  })

  it('blocks empty company name', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('Valid Name')
    cy.get('[data-cy="supplier-contact"]').type('+919999999999')
    cy.get('[data-cy="supplier-email"]').type('test@test.com')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Company name required')
  })

  it('blocks invalid company characters', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-name"]').type('Valid Name')
    cy.get('[data-cy="supplier-contact"]').type('+919999999999')
    cy.get('[data-cy="supplier-email"]').type('test@test.com')
    cy.get('[data-cy="supplier-company"]').type('@@@###')
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains('.Toastify__toast', 'Company name should contain only alphabets, numbers and spaces')
  })


  it('creates a new supplier', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()

    cy.get('[data-cy="supplier-name"]').type(supplierName)
    cy.get('[data-cy="supplier-contact"]').type('+919999999999')
    cy.get('[data-cy="supplier-email"]').type('supplier@test.com')
    cy.get('[data-cy="supplier-company"]').type('Test Company')

    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains(supplierName).should('exist')
  })


  it('edits an existing supplier', () => {
    cy.contains('td', supplierName)
      .parents('tr')
      .within(() => {
        cy.get('[data-cy="edit-supplier-btn"]').click()
      })

    cy.get('[data-cy="supplier-name"]').clear().type(updatedName)
    cy.get('[data-cy="supplier-submit"]').click()

    cy.contains(updatedName).should('exist')
  })


  it('closes modal using Cancel', () => {
    cy.get('[data-cy="add-supplier-btn"]').click()
    cy.get('[data-cy="supplier-cancel"]').click()

    cy.get('[data-cy="supplier-modal-title"]').should('not.exist')
  })


  it('deletes the same supplier', () => {
    cy.contains('td', updatedName)
      .parents('tr')
      .within(() => {
        cy.get('[data-cy="delete-supplier-btn"]').click()
      })

    cy.contains(updatedName).should('not.exist')
  })

})
