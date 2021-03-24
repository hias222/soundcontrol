import React, { Component } from 'react';
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';
import './App.css';
import { connect, sendMsg } from './api';
import Container from '@material-ui/core/Container';
import { Grid } from '@material-ui/core';
import Slider from '@material-ui/core/Slider';
import VolumeDown from '@material-ui/icons/VolumeDown';
import VolumeUp from '@material-ui/icons/VolumeUp';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: [],
      value: 1
    }
    this.handleChange = this.handleChange.bind(this)
  }

  componentDidMount() {
    connect((msg) => {
      //console.log("New Message")
      this.setState({
        chatHistory: [ msg]
      })
      //console.log(this.state);
    });
  }

  send(event) {
    if (event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }

  handleChange(event, newValue) {
    this.setState({
      value: newValue
    })

    var volValue = newValue / 100
   

    var sendValue = {
      id: 1,
      percent: volValue
    }

    //console.log(sendValue)
    sendMsg(JSON.stringify(sendValue));

  };

  render() {
    return (
      <Container maxWidth="md">
        <Header />
        <Grid container spacing={2}>
          <Grid item>
            <VolumeDown />
          </Grid>
          <Grid item xs>
            <Slider value={this.state.value} onChange={this.handleChange} aria-labelledby="continuous-slider" />
          </Grid>
          <Grid item>
            <VolumeUp />
          </Grid>
        </Grid>
        <ChatHistory chatHistory={this.state.chatHistory} />
        <ChatInput send={this.send} />
      </Container>
    );
  }
}

export default App;
