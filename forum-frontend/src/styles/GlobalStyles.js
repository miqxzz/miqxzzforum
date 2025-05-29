import { createGlobalStyle } from 'styled-components';

const GlobalStyle = createGlobalStyle`
  @import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@300;400;500;600;700&display=swap');

  * {
    font-family: 'Montserrat', Arial, sans-serif;
  }

  body {
    font-family: 'Montserrat', Arial, sans-serif;
    background-color: #e0c3fc;
    color: #333;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  h1, h2, h3, h4, h5, h6 {
    color: #a259ff;
    font-family: 'Montserrat', Arial, sans-serif;
  }

  a {
    text-decoration: none;
    color: #a259ff;
    transition: color 0.3s ease;
    font-family: 'Montserrat', Arial, sans-serif;

    &:hover {
      color: #6c2eb7;
    }
  }

  button {
    padding: 10px 15px;
    background-color: #a259ff;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    font-family: 'Montserrat', Arial, sans-serif;

    &:hover {
      background-color: #6c2eb7;
    }

    &:disabled {
      background-color: #ccc;
      cursor: not-allowed;
    }
  }

  input[type="text"],
  input[type="password"],
  textarea {
    padding: 8px;
    border: 1px solid #a259ff;
    border-radius: 4px;
    margin-bottom: 10px;
    width: 100%;
    box-sizing: border-box;
    font-family: 'Montserrat', Arial, sans-serif;
  }

  input::placeholder,
  textarea::placeholder {
    font-family: 'Montserrat', Arial, sans-serif;
    color: #bfa6e6;
    opacity: 1;
  }
`;

export default GlobalStyle;