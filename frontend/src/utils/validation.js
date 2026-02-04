import { toast } from "react-toastify";

export const validateSupplierForm = (supplier) => {
  const name = supplier.name.trim()
  const contact = supplier.contact.trim()
  const email = supplier.email.trim()
  const company = supplier.company.trim()

  if (!name) {
    toast.error("Supplier name required")
    return false
  }

  if (name.length > 50) {
    toast.error("Supplier name too long")
    return false
  }

  if (name.length < 2) {
    toast.error("Supplier name too short")
    return false
  }

  const validName = /^[A-Za-z ]+$/
  if (!validName.test(name)) {
    toast.error("Supplier name should only contain alphabets and spaces")
    return false
  }

  if (!/[A-Za-z]/.test(name)) {
    toast.error("Supplier name must contain at least one letter")
    return false
  }

  if (!contact) {
    toast.error("Contact number required")
    return false
  }

  const validPhone = /^[0-9+\-() ]+$/
  if (!validPhone.test(contact)) {
    toast.error("Contact should contain only numbers, +, -, () or spaces")
    return false
  }

  if (contact.length < 7 || contact.length > 15) {
  toast.error("Contact number must be between 7 and 15 characters")
  return false
}

  if (!email) {
    toast.error("Email is required")
    return false
  }

  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(email)) {
    toast.error("Invalid email format (e.g., user@example.com)")
    return false
  }

  if (!company) {
    toast.error("Company name required")
    return false
  }

  if (company.length > 50) {
    toast.error("Company name too long")
    return false
  }

  if (company.length < 2) {
    toast.error("Company name too short")
    return false
  }

  const validCompany = /^[A-Za-z0-9 ]+$/
  if (!validCompany.test(company)) {
    toast.error("Company name should contain only alphabets, numbers and spaces")
    return false
  }

  if (!/[A-Za-z]/.test(company)) {
    toast.error("Company name must contain at least one letter")
    return false
  }

  return true
}
