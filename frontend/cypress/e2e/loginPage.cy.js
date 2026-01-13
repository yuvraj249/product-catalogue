/* eslint-disable no-unused-expressions */
/* eslint-disable no-undef */

describe('Login Page - Auth and Validation', () => {

  beforeEach(() => {
    cy.clearLocalStorage()
    cy.visit('/')
  })

  it('loads Login Page UI', () => {
    cy.contains('Sign in to your account').should('be.visible')
    cy.get('input[placeholder="abc@gmail.com"]').should('be.visible')
    cy.get('input[placeholder="password"]').should('be.visible')
    cy.contains('Sign In').should('be.visible')
  })


  it('blocks empty email', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="password"]').type('Valid@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })



  it('blocks invalid email format', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('invalidemail.com')
    cy.get('input[placeholder="password"]').type('Valid@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks email starting with @', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('@test.com')
    cy.get('input[placeholder="password"]').type('Valid@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks email with consecutive dots', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('test..mail@gmail.com')
    cy.get('input[placeholder="password"]').type('Valid@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks small passwords', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('Valid')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks password without uppercase', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('valid@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks password without lowercase', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('VALID@123')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks password without number', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('VALID@valid')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })


  it('blocks password without special character', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('Valid12345')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 400)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Login failed')
  })



  it('blocks invalid credentials', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('tester@tester.com')
    cy.get('input[placeholder="password"]').type('Valid@12345')
    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 401)

    cy.get('.Toastify__toast')
      .should('be.visible')
      .and('contain.text', 'Incorrect email or password')
  })

  

  it('valid log in with system_admin credentials', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('input[placeholder="abc@gmail.com"]').type('yuvrajbisht41@gmail.com')
    cy.get('input[placeholder="password"]').type('Yuvraj@2411')

    cy.contains('Sign In').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 200)

    cy.url({ timeout: 15000 }).should('include', '/admin/categories')

    cy.window().then(win => {
      expect(win.localStorage.getItem('auth_token')).to.exist
    })
  })


  it('blocks unauthenticated user from admin routes', () => {
    cy.clearLocalStorage()
    cy.visit('/admin/categories')

    cy.url().should('eq', Cypress.config().baseUrl + '/')
  })

})
