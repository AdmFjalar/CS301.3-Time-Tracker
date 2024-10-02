import React from "react";
import Classprops from "./Classprops";
import Functionprops from './Functionprops';


class Timetracker extends React.Component {
    render(){
        return (
            <div>
                <Classprops name="Manager">
                    <p>Child Component</p>
                    </Classprops>
                <Classprops name="Employer">
                    <button>Click</button>
                    </Classprops>
                <Classprops name="Employee"/>
                <Functionprops name="Admin"></Functionprops>
            </div>
        );
    }
}


export default Timetracker;