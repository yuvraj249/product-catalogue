import styled from "styled-components";
import ReactSelect from "react-select";

export const SupplierSelect = styled(ReactSelect)`
  .react-select__control {
    border-radius: 8px;
    background: #f3f8fb;
    border: 1px solid rgba(42,123,155,0.30);
    min-height: 38px;
  }

  .react-select__control--is-focused {
    border-color: #2A7B9B;
    box-shadow: 0 0 0 3px rgba(42,123,155,0.15);
    background: white;
  }

  .react-select__menu-list {
    max-height: 180px;
    overflow-y: auto;
  }

  .react-select__option {
    padding: 6px 10px;
    font-size: 14px;
  }

  .react-select__option--is-selected {
    background: rgba(42,123,155,0.12);
    color: #2A7B9B;
  }

  .react-select__option--is-focused {
    background: rgba(42,123,155,0.08);
  }

  .react-select__value-container {
    padding: 4px 8px;
  }

  .react-select__indicator-separator {
    display: none;
  }
`;
