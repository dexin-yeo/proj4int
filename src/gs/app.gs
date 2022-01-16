function insert() {
  var ui = SpreadsheetApp.getUi();
  // get current sheet
  var activeSheet = SpreadsheetApp.getActiveSpreadsheet();

  // prompting user for spreadsheet name to edit
  var prompt = ui.prompt(
    'Enter the spreadsheet name to insert row',
    ui.ButtonSet.OK_CANCEL);

  // response text
  var res = prompt.getResponseText();
  // response button
  var button = prompt.getSelectedButton();

  if (button == ui.Button.OK) {
    var toEditSheet = activeSheet.getSheetByName(res);
    if (toEditSheet == null) {
      // console.log(":")
      ui.alert('Unable to find spreadsheet: ' + res)
      return;
    }

    // go to the sheet to be edited
    activeSheet.setActiveSheet(toEditSheet);

    prompt = ui.prompt(
      'Enter the data to insert by the order: table_name,column_name,ordinal_position,column_default,is_nullable,data_type,character_maximum_length',
      'e.g. my_roads,tmpcol1,20,,NO,integer,',
      ui.ButtonSet.OK_CANCEL);

    // response text
    res = prompt.getResponseText();
    // response button
    button = prompt.getSelectedButton();

    if (button == ui.Button.OK) {
      var arr = res.split(",");
      // get ss data
      var data = activeSheet.getDataRange().getValues();

      // get the cols name
      var cols = data[0];
      
      // remove all the special stuff
      for (let i = 0; i < cols.length; i++) {
        cols[i] = cols[i].slice(3,cols[i].length).replace(/ /g,'');
      }

      // creates array for values
      var values = [new Array(11)];
      var col = cols.indexOf("table_name");
      values[0][col] = arr[0];

      // get the empty cell row base on table_name entry
      var cr = 0;
      while ( data[cr] && data[cr][col] != "" ) {
        cr++;
      }
      cr+=1;

      col = cols.indexOf('column_name');
      console.log(col)
      values[0][col] = arr[1];
      col = cols.indexOf("ordinal_position");
      values[0][col] = arr[2];
      col = cols.indexOf("column_default");
      values[0][col] = arr[3];
      col = cols.indexOf("is_nullable");
      values[0][col] = arr[4];
      col = cols.indexOf("data_type");
      values[0][col] = arr[5];
      col = cols.indexOf("character_maximum_length");
      values[0][col] = arr[6];

      // // get current date
      // var currDate = new Date();
      // var spreadsheet = SpreadsheetApp.getActive();
      // var timezone = spreadsheet.getSpreadsheetTimeZone(); 
      // var format = 'HH:mm:ss';
      // var time = Utilities.formatDate(currDate, timezone, format); 
      // format = "yyyy-MM-dd"
      // var date = Utilities.formatDate(currDate, timezone, format); 
      // // spr.getRange("D2").setFormula('=datevalue("'+ date +'")+timevalue("'+time+'")')
      // col = cols.indexOf("Createdtime");
      // values[0][col] = '=datevalue("' + date + '")+timevalue("' + time+ '")';

      // col = cols.indexOf("Lasteditedtime");
      // values[0][col] = '=datevalue("' + date + '")+timevalue("' + time+ '")';


      // row#, col#, #ofRows, #ofCols
      var cell = activeSheet.getActiveSheet().getRange(cr,1,1,cols.length);
      cell.setValues(values);
    }
  }
}

function update() {
  var ui = SpreadsheetApp.getUi();
  // get current sheet
  var activeSheet = SpreadsheetApp.getActiveSpreadsheet();

  // prompting user for spreadsheet name to edit
  var prompt = ui.prompt(
    'Enter the spreadsheet name to update row',
    ui.ButtonSet.OK_CANCEL);

  // response text
  var res = prompt.getResponseText();
  // response button
  var button = prompt.getSelectedButton();

  if (button == ui.Button.OK) {
    var toEditSheet = activeSheet.getSheetByName(res);
    if (toEditSheet == null) {
      // console.log(":")
      ui.alert('Unable to find spreadsheet: ' + res);
      return;
    }

    // go to the sheet to be edited
    activeSheet.setActiveSheet(toEditSheet);

    prompt = ui.prompt(
      'Search by row and col?',
      'e.g. 2,3',
      ui.ButtonSet.OK_CANCEL);

    res = prompt.getResponseText();
    button = prompt.getSelectedButton();

    if (button == ui.Button.OK) {
      var arr = res.split(",");
      prompt = ui.prompt(
        'What value to edit to?',
        'e.g. edited_name',
        ui.ButtonSet.OK_CANCEL);
        res = prompt.getResponseText();
        button = prompt.getSelectedButton();
        if (button == ui.Button.OK) {
          // row#, col#, #ofRows, #ofCols
          var cell = activeSheet.getActiveSheet().getRange(arr[0],arr[1]);
          cell.setValues([[res]]);
        }
    } else if (button == ui.Button.CANCEL) {
      prompt = ui.prompt(
        'Search by page_id?',
        'e.g. 015eab9e-3689-491c-9775-20dc74a9d06d',
        ui.ButtonSet.OK_CANCEL);

      var id = prompt.getResponseText();
      button = prompt.getSelectedButton();
      if (button == ui.Button.OK) {
        // get ss data
        var data = activeSheet.getDataRange().getValues();

        // get the cols name
        var cols = data[0];
        
        // remove all the special stuff
        for (let i = 0; i < cols.length; i++) {
          cols[i] = cols[i].slice(3,cols[i].length).replace(/ /g,'');
        }

        // creates array for values
        var values = [new Array(11)];
        var col = cols.indexOf("PageID");
        var row = 1;
        // get row
        for (let i = 0; i < data.length; i++) {
          if (id == data[i][col]) {
            row = i;
            for (let j = 0; j < cols.length; j++) {
              values[0][j] = data[i][j];
            }
            break
          }
        }

        while (button == ui.Button.OK) {
          prompt = ui.prompt(
            'What column to edit?',
            'e.g. column_name',
            ui.ButtonSet.OK_CANCEL);
          res = prompt.getResponseText();
          button = prompt.getSelectedButton();
          if (button == ui.Button.OK) {
            col = cols.indexOf(res);
            if (col < 0) {
              ui.alert('Unable to find column: ' + res);
              return
            }
            prompt = ui.prompt(
              'What value to change to?',
              'e.g. edited_name',
              ui.ButtonSet.OK_CANCEL);
            res = prompt.getResponseText();
            button = prompt.getSelectedButton();
            if (button == ui.Button.OK) {
              values[0][col] = res;
            }
          }
        }
        // row#, col#, #ofRows, #ofCols
        var cell = activeSheet.getActiveSheet().getRange(row+1,1,1,cols.length);
        cell.setValues(values);
      }
    }
  }
}