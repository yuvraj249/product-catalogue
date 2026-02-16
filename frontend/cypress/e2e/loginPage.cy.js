/* eslint-disable no-unused-expressions */
/* eslint-disable no-undef */

describe('Login Page - Auth and Validation', () => {

  beforeEach(() => {
    cy.clearLocalStorage()
    cy.visit('/')
  })

  it('loads Login Page UI', () => {
    cy.contains('Sign in to your account').should('be.visible')
    cy.get('[data-cy="email-input"]').should('be.visible')
    cy.get('[data-cy="password-input"]').should('be.visible')
    cy.get('[data-cy="login-btn"]').should('be.visible')
  })


  it('blocks empty email', () => {
    cy.get('[data-cy="password-input"]').type('Valid@123')
    cy.get('[data-cy="login-btn"]').click()

    cy.contains('.Toastify__toast', 'Email is required')
  })

  it('blocks invalid email format', () => {
    cy.get('[data-cy="email-input"]').type('invalidemail.com')
    cy.get('[data-cy="password-input"]').type('Valid@123')
    cy.get('[data-cy="login-btn"]').click()

    cy.contains('.Toastify__toast', 'Enter a valid email')
  })

  it('blocks email starting with @', () => {
    cy.get('[data-cy="email-input"]').type('@test.com')
    cy.get('[data-cy="password-input"]').type('Valid@123')
    cy.get('[data-cy="login-btn"]').click()

    cy.contains('.Toastify__toast', 'Enter a valid email')
  })

  it('blocks email with consecutive dots', () => {
    cy.get('[data-cy="email-input"]').type('test..mail@gmail.com')
    cy.get('[data-cy="password-input"]').type('Valid@123')
    cy.get('[data-cy="login-btn"]').click()

    cy.contains('.Toastify__toast', 'Login failed')
  })

  it('blocks empty password', () => {
    cy.get('[data-cy="email-input"]').type('tester@tester.com')
    cy.get('[data-cy="login-btn"]').click()

    cy.contains('.Toastify__toast', 'Password is required')
  })


  it('blocks invalid credentials', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('[data-cy="email-input"]').type('tester@tester.com')
    cy.get('[data-cy="password-input"]').type('Wrong@123')
    cy.get('[data-cy="login-btn"]').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 401)

    cy.contains('.Toastify__toast', 'Incorrect email or password')
  })


  it('valid login with system_admin credentials', () => {
    cy.intercept('POST', '**/auth/login').as('loginRequest')

    cy.get('[data-cy="email-input"]').type('yuvrajbisht41@gmail.com')
    cy.get('[data-cy="password-input"]').type('Yuvraj@2411')
    cy.get('[data-cy="login-btn"]').click()

    cy.wait('@loginRequest')
      .its('response.statusCode')
      .should('eq', 200)


    cy.url().should('include', '/admin/categories')

    cy.window().then(win => {
      expect(win.localStorage.getItem('auth_token')).to.exist
    })
  })

})
