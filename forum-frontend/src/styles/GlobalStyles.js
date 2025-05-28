import { createGlobalStyle } from 'styled-components';

const GlobalStyle = createGlobalStyle`
  @import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@300;400;500;600;700&display=swap');

  * {
    font-family: 'Montserrat', sans-serif;
  }

  body {
    font-family: 'Montserrat', sans-serif;
    background-color: #f4f4f4;
    color: #333;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  h1, h2, h3, h4, h5, h6 {
    color: #8e44ad;
    font-weight: 600;
  }

  a {
    text-decoration: none;
    color: #9b59b6;
    transition: color 0.3s ease;

    &:hover {
      color: #8e44ad;
    }
  }

  button {
    padding: 10px 15px;
    background-color: #9b59b6;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    font-weight: 500;

    &:hover {
      background-color: #8e44ad;
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
    border: 1px solid #ddd;
    border-radius: 4px;
    margin-bottom: 10px;
    width: 100%;
    box-sizing: border-box;
    font-family: 'Montserrat', sans-serif;
  }
`;

export default GlobalStyle;