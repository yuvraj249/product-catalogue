import { FormGroup, Label } from "./Styles";

const FormField = ({ label, required = false, children }) => {
  return (
    <FormGroup>
      <Label>
        {label} {required && "*"}
      </Label>
      {children}
    </FormGroup>
  );
};

export default FormField;
