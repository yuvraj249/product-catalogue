/* eslint-disable no-undef */

describe('Categories Page - System Admin', () => {

  const categoryName = "Test category"
  const updatedName = "Update category"

  beforeEach(() => {
    cy.loginAsSystemAdmin()
    cy.visit('/admin/categories')
  })

  it('loads categories page', () => {
    cy.contains('Categories').should('be.visible')
    cy.get('[data-cy="add-category-btn"]').should('be.visible')
  })

  it('opens Add Category modal', () => {
    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-modal-title"]')
      .should('contain.text', 'Add Category')
  })

  it('blocks empty category name', () => {
    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-submit"]').click()

    cy.contains('.Toastify__toast', 'Category name is required')
  })


  it('blocks name shorter than 2 characters', () => {
    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-name"]').type('A')
    cy.get('[data-cy="category-submit"]').click()

    cy.contains('.Toastify__toast', 'Category name must be at least 2 characters')
  })


  it('blocks name longer than 50 characters', () => {
    const longName = 'F'.repeat(60)

    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-name"]').type(longName)
    cy.get('[data-cy="category-submit"]').click()

    cy.contains('.Toastify__toast', 'Category name must be under 50 characters')
  })


  it('blocks description longer than 200 characters', () => {
    const longDesc = 'D'.repeat(250)

    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-name"]').type('Valid Name')
    cy.get('[data-cy="category-description"]').type(longDesc)
    cy.get('[data-cy="category-submit"]').click()

    cy.contains('.Toastify__toast', 'Description must be under 200 characters')
  })


  it('creates a new category', () => {
    cy.get('[data-cy="add-category-btn"]').click()

    cy.get('[data-cy="category-name"]').type(categoryName)
    cy.get('[data-cy="category-description"]').type('Test Description')
    cy.get('[data-cy="category-submit"]').click()

    cy.contains(categoryName).should('exist')
  })


  it('blocks duplicate category name', () => {
    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-name"]').type(categoryName)
    cy.get('[data-cy="category-submit"]').click()

    cy.contains('.Toastify__toast', 'Category with this name already exists')
  })


  it('edits the same category', () => {
    cy.contains('td', categoryName)
      .parents('tr')
      .within(() => {
        cy.get('[data-cy="edit-category-btn"]').click()
      })

    cy.get('[data-cy="category-name"]')
      .clear()
      .type(updatedName)

    cy.get('[data-cy="category-submit"]').click()

    cy.contains(updatedName).should('exist')
  })

  it('closes modal using Cancel button', () => {
    cy.get('[data-cy="add-category-btn"]').click()
    cy.get('[data-cy="category-cancel"]').click()

    cy.get('[data-cy="category-modal-title"]').should('not.exist')
  })


  it('deletes the same category', () => {
    cy.contains('td', updatedName)
      .parents('tr')
      .within(() => {
        cy.get('[data-cy="delete-category-btn"]').click()
      })

    cy.contains(updatedName).should('not.exist')
  })

})
