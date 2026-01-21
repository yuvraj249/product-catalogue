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
    email: 'admin1@test.com',
    password: 'Admin@123'
  }).then((res) => {
    window.localStorage.setItem('auth_token', res.body.token)
  })
})  

Cypress.Commands.add('seedCategories', () => {
  const categories = [
    { name: 'Electronics', desc: 'Electronic items' },
    { name: 'Groceries', desc: 'Daily groceries' }
  ]

  cy.wrap([]).as('categoryIds')

  categories.forEach(cat => {
    cy.request({
      method: 'POST',
      url: 'http://localhost:8080/categories',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('auth_token')}`
      },
      body: {
        category_name: cat.name,
        category_description: cat.desc
      }
    }).then(res => {
      cy.get('@categoryIds').then(ids => {
        ids.push(res.body.category_id)
        cy.wrap(ids).as('categoryIds')
      })
    })
  })
})


Cypress.Commands.add('seedSuppliers', () => {
  const suppliers = [
    { name: 'Supplier A', contact: '+911111111111', email: 's1@test.com', company: 'Company A' },
    { name: 'Supplier B', contact: '+922222222222', email: 's2@test.com', company: 'Company B' },
    { name: 'Supplier C', contact: '+933333333333', email: 's3@test.com', company: 'Company A' },
    { name: 'Supplier D', contact: '+944444444444', email: 's4@test.com', company: 'Company A' },
    { name: 'Supplier E', contact: '+955555555555', email: 's5@test.com', company: 'Company E' },
  ]

  suppliers.forEach(s => {
    cy.request({
      method: 'POST',
      url: 'http://localhost:8080/suppliers',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('auth_token')}`
      },
      body: {
        name: s.name,
        contact_info: s.contact,
        email: s.email,
        company: s.company
      },
      failOnStatusCode: true    })
  })
})


Cypress.Commands.add('seedUsers', () => {
  cy.request({
    method: 'GET',
    url: 'http://localhost:8080/suppliers',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('auth_token')}`
    }
  }).then((res) => {
    const suppliers = res.body.suppliers

    const users = [
      { name: 'Supplier Admin 1', email: 'admin1@test.com', password: 'Admin@123', supplier_id: suppliers[0].supplier_id },
      { name: 'Supplier Admin 2', email: 'admin2@test.com', password: 'Admin@123', supplier_id: suppliers[1].supplier_id },
      { name: 'Supplier Admin 3', email: 'admin3@test.com', password: 'Admin@123', supplier_id: suppliers[2].supplier_id },
      { name: 'Supplier Admin 4', email: 'admin4@test.com', password: 'Admin@123', supplier_id: suppliers[3].supplier_id },
      { name: 'Supplier Admin 5', email: 'admin5@test.com', password: 'Admin@123', supplier_id: suppliers[4].supplier_id },
    ]

    users.forEach(user => {
      cy.request({
        method: 'POST',
        url: 'http://localhost:8080/users/supplier-admin',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: user,
        failOnStatusCode: true      })
    })
  })
})



Cypress.Commands.add('seedProducts', () => {
  cy.get('@categoryIds').then(ids => {

    cy.wrap([]).as('productIds')

    const products = [
      { name: 'Laptop', cost: 50000, cat: ids[0], desc: 'Gaming Laptop' },
      { name: 'Rice Bag', cost: 1200, cat: ids[1], desc: '25kg Rice' },
    ]

    products.forEach(p => {
      cy.request({
        method: 'POST',
        url: 'http://localhost:8080/products',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: {
          product_name: p.name,
          product_description: p.desc,
          product_cost: p.cost,
          product_category_id: p.cat,
          discount_type: "",
          discount_value: 0
        }
      }).then(res => {
        cy.get('@productIds').then(prodIds => {
          prodIds.push(res.body.product_id)
          cy.wrap(prodIds).as('productIds')
        })
      })
    })
  })
})


Cypress.Commands.add('seedStock', () => {
  cy.get('@productIds').then(ids => {

    const stock = [
      { product_id: ids[0], quantity: 5, movement_type: 'IN',  reason: 'Initial stock' },
      { product_id: ids[1], quantity: 10, movement_type: 'IN', reason: 'Warehouse refill' },
    ]

    stock.forEach(s => {
      cy.request({
        method: 'POST',
        url: 'http://localhost:8080/stock_movements',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: s,
        failOnStatusCode: true
      })
    })
  })
})

Cypress.Commands.add('cleanupAll', () => {
  cy.request({
    method: 'DELETE',
    url: 'http://localhost:8080/test/cleanup/all',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('auth_token')}`
    }
  })
})




