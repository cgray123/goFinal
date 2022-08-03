import React, { Component } from 'react';
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';
import NameInput from './components/NameInput/NameInput';
import './App.css';
import { connect, sendMsg } from './api';



class App extends Component {
  
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: []
    }
   
  }
 
  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...prevState.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }
/* send input to be handle as a msg */
  send(event) {
    if (event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }
  /* send input to set username */
  name(event){
    if (event.keyCode === 13) {
      sendMsg("~"+event.target.value);
      event.target.value = "";
    }
  }

 
  render() {
    return (
      <div className="App">
        <Header />
        <ChatHistory chatHistory={this.state.chatHistory}/>
        <NameInput name={this.name}/>
       <ChatInput send={this.send} />
        
      </div>
    );
  }
}

export default App;
