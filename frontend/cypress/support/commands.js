/* eslint-disable no-undef */

Cypress.Commands.add('loginAsSystemAdmin', () => {
  cy.request('POST', 'http://localhost:8080/auth/login', {
    email: 'yuvrajbisht41@gmail.com',
    password: 'Yuvraj@2411'
  }).then((res) => {
    window.localStorage.setItem('auth_token', res.body.token)
  })
})

Cypress.Commands.add('loginAsSupplierAdmin', () => {
  cy.request('POST', 'http://localhost:8080/auth/login', {
    email: 'tester@tester.com',
    password: 'Tester@1234'
  }).then((res) => {
    window.localStorage.setItem('auth_token', res.body.token)
  })
})