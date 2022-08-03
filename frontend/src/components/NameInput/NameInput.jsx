import React, { Component } from 'react';
import './NameInput.scss';

class  NameInput extends Component {
  state = {
    isVis:true
  }

  toggle = ()=>{
    this.setState({
      isVis:false
    })
  }
  render() {
    
    return (
      <div className=' NameInput'>
       { /* when enter is hit, makes nameinput hidden by calling toggle */}
        { this.state.isVis ? <input onKeyDown={this.props.name} onKeyPress={(e) => e.key === 'Enter' && this.toggle()} placeholder="Enter your Name before sending a Message"/>:null}
      </div>
    );
  };

}

export default  NameInput;