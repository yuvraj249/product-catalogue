/* eslint-disable no-undef */

describe('Products Page - Supplier Admin', () => {

  const productName = "Test Product"
  const updatedName = "Updated Product"

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
    cy.loginAsSupplierAdmin()
    cy.visit('/admin/products')
  })

  after(() => {
    cy.loginAsSystemAdmin()
    cy.cleanupAll()
  })

  it('loads Products page', () => {
    cy.get('[data-cy="products-title"]').should('contain', 'Products')
    cy.get('[data-cy="add-product-btn"]').should('be.visible')
    cy.get('[data-cy="data-table"]').should('be.visible')
  })

  it('opens Add Product modal', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-modal-title"]').should('contain', 'Add Product')
  })


  it('blocks empty form submission', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('Product name required').should('exist')
  })

  it('blocks name with only special characters', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type('@@@')
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('only alphabets').should('exist')
  })

  it('blocks name without letters', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type('12345')
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('at least one letter').should('exist')
  })

  it('blocks too short name', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type('A')
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('too short').should('exist')
  })

  it('blocks too long name', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type('A'.repeat(60))
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('too long').should('exist')
  })

  it('blocks missing category', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type(productName)
    cy.get('[data-cy="product-cost"]').type('500')
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('select a valid category').should('exist')
  })

  it('blocks invalid cost', () => {
    cy.get('[data-cy="add-product-btn"]').click()
    cy.get('[data-cy="product-name"]').type(productName)
    cy.get('[data-cy="product-cost"]').type('-10')
    cy.get('[data-cy="product-submit-btn"]').click()
    cy.contains('realistic product cost').should('exist')
  })

  it('blocks discount without type', () => {
  cy.get('[data-cy="add-product-btn"]').click()

  cy.get('[data-cy="product-name"]').type(productName)
  cy.get('[data-cy="product-cost"]').type('500')
  cy.get('[data-cy="product-discount-value"]').type('50')
  cy.get('[data-cy="product-submit-btn"]').click()
  cy.contains(productName).should('not.exist')
})

  it('blocks percent discount > 100', () => {
    cy.get('[data-cy="add-product-btn"]').click()

    cy.get('[data-cy="product-name"]').type(productName)
    cy.get('[data-cy="product-cost"]').type('500')
    cy.get('[data-cy="product-discount-type"]').type('percent')
    cy.get('[data-cy="product-discount-value"]').type('150')

    cy.get('[data-cy="product-submit-btn"]').click()

    cy.contains(productName).should('not.exist')
  })

  it('blocks flat discount > cost', () => {
    cy.get('[data-cy="add-product-btn"]').click()

    cy.get('[data-cy="product-name"]').type(productName)
    cy.get('[data-cy="product-cost"]').type('500')
    cy.get('[data-cy="product-discount-type"]').type('flat')
    cy.get('[data-cy="product-discount-value"]').type('800')

    cy.get('[data-cy="product-submit-btn"]').click()

    cy.contains(productName).should('not.exist')
  })


  it('creates a product', () => {
    cy.get('[data-cy="add-product-btn"]').click()

    cy.get('[data-cy="product-name"]').type(productName)

    cy.get('.react-select__control').click()
    cy.get('.react-select__menu').contains('Electronics').click()

    cy.get('[data-cy="product-cost"]').type('500')

    cy.get('[data-cy="product-submit-btn"]').click()

    cy.contains(productName, { timeout: 8000 }).should('exist')
  })

  it('edits a product', () => {
    cy.contains('td', productName)
      .parents('tr')
      .find('[data-cy="edit-product-btn"]')
      .click()

    cy.get('[data-cy="product-name"]').clear().type(updatedName)
    cy.get('[data-cy="product-submit-btn"]').click()

    cy.contains(updatedName).should('exist')
  })

  it('searches products', () => {
    cy.get('[data-cy="product-search"]').type(updatedName)
    cy.contains(updatedName).should('exist')
  })

  it('deletes a product', () => {
    cy.contains('td', updatedName)
      .parents('tr')
      .find('[data-cy="delete-product-btn"]')
      .click()

    cy.contains(updatedName).should('not.exist')
  })

})
