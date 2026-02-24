import styled from "styled-components";
import Select from "react-select";

export const SelectBox = styled(Select)`
  .react-select__control {
    min-height: 52px;
    border-radius: 12px;
    border: 1px solid rgba(42, 123, 155, 0.3);
    background: #f3f8fb;
    padding: 6px 8px;
    box-shadow: none;
  }

  .react-select__control--is-focused {
    border-color: #2a7b9b;
    box-shadow: 0 0 0 3px rgba(42, 123, 155, 0.15);
    background: white;
  }

  .react-select__option--is-selected {
    background: #2a7b9b;
    color: white;
  }

  .react-select__option--is-focused {
    background: rgba(42, 123, 155, 0.12);
  }
`;
