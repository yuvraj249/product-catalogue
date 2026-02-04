import { toast } from "react-toastify";
import api from "../Api/axios";


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


export const validateUserForm = (userObj, updatingId, form) => {
  const name = userObj.name.trim()
  const email = userObj.email.trim()
  const password = userObj.password
  const supplierId = userObj.supplier_id

  if (!name) {
    toast.error("User name is required")
    return false
  }

  if (name.length < 2) {
    toast.error("User name must be at least 2 characters")
    return false
  }

  if (!email) {
    toast.error("Email cannot be empty")
    return false
  }

  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(email)) {
    toast.error("Invalid email format (e.g., user@example.com)")
    return false
  }

  if (!updatingId) {
    if (!password) {
      toast.error("Password is required")
      return false
    }

    if (password.length < 8) {
      toast.error("Password must be at least 8 characters long")
      return false
    }

    if (!/^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9])(?=.*[!@#$%^&*(),.?":{}|<>_\-+=~`])/.test(password)) {
          toast.error("Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
          return false
    }

  }

  if (form.role !== "system_admin") {
    if (!supplierId || Number(supplierId) <= 0) {
      toast.error("Please select a valid supplier")
      return false
    }
  }

  return true
}

export const validateProductForm = (product) => {
    const name = product.product_name.trim();
    const desc = product.product_description?.trim() || "";
    const cost = Number(product.product_cost);
    const categoryId = Number(product.product_category_id);
    const discountType = product.discount_type.trim().toLowerCase();
    const discountValue = Number(product.discount_value);

    if (!name) {
      toast.error("Product name required");
      return false;
    }

    if (name.length > 50) {
      toast.error("Product name too long");
      return false;
    }

    if (name.length < 2) {
      toast.error("Product name too short");
      return false;
    }

    const validName = /^[A-Za-z0-9 ]+$/;
    if (!validName.test(name)) {
      toast.error("Product name should contain only alphabets, numbers and spaces");
      return false;
    }

    if (!/[A-Za-z]/.test(name)) {
      toast.error("Product name must contain at least one letter");
      return false;
    }

    if (cost <= 0 || cost > 9999999999.5) {
      toast.error("Please enter a realistic product cost");
      return false;
    }

    if (desc && desc.length > 1500) {
      toast.error("Description too long");
      return false;
    }

    if (!categoryId || categoryId <= 0) {
      toast.error("Please select a valid category");
      return false;
    }

    if (discountValue !== 0 && !discountType) {
      toast.error("Discount type required if discount value is provided");
      return false;
    }

    if (discountType) {
      if (discountType !== "flat" && discountType !== "percent") {
        toast.error("Discount type must be flat or percent");
        return false;
      }

      if (discountValue === 0) {
        toast.error("Discount value required when discount type is selected");
        return false;
      }

      if (discountValue < 0) {
        toast.error("Discount value must be >= 0");
        return false;
      }

      if (discountType === "percent" && discountValue > 100) {
        toast.error("Percent discount cannot exceed 100");
        return false;
      }

      if (discountType === "flat" && discountValue > cost) {
        toast.error("Flat discount cannot exceed product cost");
        return false;
      }
    }

    return true;
  }

  
export const validateCategoryForm = async (categoryObj, updatingId) => {
  const name = categoryObj.name.trim()
  const desc = categoryObj.description.trim()

  if (!name) {
    toast.error("Category name is required")
    return false
  }

  if (name.length < 2) {
    toast.error("Category name must be at least 2 characters")
    return false
  }

  if (name.length > 50) {
    toast.error("Category name must be under 50 characters")
    return false
  }

  if (desc.length > 200) {
    toast.error("Description must be under 200 characters")
    return false
  }

  try {
    const res = await api.get("/categories", {
      params: { q: name }
    })

    const exists = res.data.categories?.some(
      c => c.category_name.toLowerCase() === name.toLowerCase()
    )

    if (exists && !updatingId) {
      toast.error("Category with this name already exists")
      return false
    }
  } catch {
    toast.error("Failed to validate category name")
    return false
  }

  return true
}