package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

Dowloaded from: http://download.shredderchess.com/div/uci.zip

Description of the universal chess interface (UCI) April 2006
=============================================================

* The specification is independent of the operating system. For Windows,
  the engine is a normal exe file, either a console or "real" windows application.

* all communication is done via standard input and output with text commands,

* The engine should boot and wait for input from the GUI,
  the engine should wait for the "isready" or "setoption" command to set up its internal parameters
  as the boot process should be as quick as possible.

* the engine must always be able to process input from stdin, even while thinking.

* all command strings the engine receives will end with '\n',
  also all commands the GUI receives should end with '\n',
  Note: '\n' can be 0x0d or 0x0a0d or any combination depending on your OS.
  If you use Engine and GUI in the same OS this should be no problem if you communicate in text mode,
  but be aware of this when for example running a Linux engine in a Windows GUI.

* arbitrary white space between tokens is allowed
  Example: "debug on\n" and  "   debug     on  \n" and "\t  debug \t  \t\ton\t  \n"
  all set the debug mode of the engine on.

* The engine will always be in forced mode which means it should never start calculating
  or pondering without receiving a "go" command first.

* Before the engine is asked to search on a position, there will always be a position command
  to tell the engine about the current position.

* by default all the opening book handling is done by the GUI,
  but there is an option for the engine to use its own book ("OwnBook" option, see below)

* if the engine or the GUI receives an unknown command or token it should just ignore it and try to
  parse the rest of the string in this line.
  Examples: "joho debug on\n" should switch the debug mode on given that joho is not defined,
            "debug joho on\n" will be undefined however.

* if the engine receives a command which is not supposed to come, for example "stop" when the engine is
  not calculating, it should also just ignore it.


Move format:
------------

The move format is in long algebraic notation.
A nullmove from the Engine to the GUI should be sent as 0000.
Examples:  e2e4, e7e5, e1g1 (white short castling), e7e8q (for promotion)


Examples:
---------

This is how the communication when the engine boots can look like:

GUI     engine

// tell the engine to switch to UCI mode
uci

// engine identify
      id name Shredder
		id author Stefan MK

// engine sends the options it can change
// the engine can change the hash size from 1 to 128 MB
		option name Hash type spin default 1 min 1 max 128

// the engine supports Nalimov endgame tablebases
		option name NalimovPath type string default <empty>
		option name NalimovCache type spin default 1 min 1 max 32

// the engine can switch off Nullmove and set the playing style
	   option name Nullmove type check default true
  		option name Style type combo default Normal var Solid var Normal var Risky

// the engine has sent all parameters and is ready
		uciok

// Note: here the GUI can already send a "quit" command if it just wants to find out
//       details about the engine, so the engine should not initialize its internal
//       parameters before here.
// now the GUI sets some values in the engine
// set hash to 32 MB
setoption name Hash value 32

// init tbs
setoption name NalimovCache value 1
setoption name NalimovPath value d:\tb;c\tb

// waiting for the engine to finish initializing
// this command and the answer is required here!
isready

// engine has finished setting up the internal values
		readyok

// now we are ready to go

// if the GUI is supporting it, tell the engine that is is
// searching on a game that it hasn't searched on before
ucinewgame

// if the engine supports the "UCI_AnalyseMode" option and the next search is supposed to
// be an analysis, the GUI should set "UCI_AnalyseMode" to true if it is currently
// set to false with this engine
setoption name UCI_AnalyseMode value true

// tell the engine to search infinite from the start position after 1.e4 e5
position startpos moves e2e4 e7e5
go infinite

// the engine starts sending infos about the search to the GUI
// (only some examples are given)


		info depth 1 seldepth 0
		info score cp 13  depth 1 nodes 13 time 15 pv f1b5
		info depth 2 seldepth 2
		info nps 15937
		info score cp 14  depth 2 nodes 255 time 15 pv f1c4 f8c5
		info depth 2 seldepth 7 nodes 255
		info depth 3 seldepth 7
		info nps 26437
		info score cp 20  depth 3 nodes 423 time 15 pv f1c4 g8f6 b1c3
		info nps 41562
		....


// here the user has seen enough and asks to stop the searching
stop

// the engine has finished searching and is sending the bestmove command
// which is needed for every "go" command sent to tell the GUI
// that the engine is ready again
		bestmove g1f3 ponder d8f6


*/

// --------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UCI Commands ----------------------------------------------
// --------------------------------------------------------------------------------------------------------------------

// starts the input loop listening to commands from the GUI
func (pos *Position) startUCIInputLoop() {

	// set up the input reader
	inputReader := bufio.NewReader(os.Stdin)

	// set up the error log buffer
	errorLogPositionBuffer := ""

	// start the command loop
	runCommandLoop := true
	for runCommandLoop {

		// get the next input command to the engine
		command, _ := inputReader.ReadString('\n')
		command = strings.TrimSpace(command)

		// respond to the command
		if command == "uci" { // remember this vs ucinewgame both have "uci" prefix
			pos.command_uci()

		} else if strings.HasPrefix(command, "debug") {
			pos.command_debug()

		} else if strings.HasPrefix(command, "isready") {
			pos.command_isReady()

		} else if strings.HasPrefix(command, "setoption") {
			pos.command_setOption()

		} else if strings.HasPrefix(command, "register") {
			pos.command_register()

		} else if strings.HasPrefix(command, "ucinewgame") {
			pos.command_uciNewGame()

		} else if strings.HasPrefix(command, "position") {
			errorLogPositionBuffer = command
			pos.command_position(command)

		} else if strings.HasPrefix(command, "go") {
			response, success := pos.command_go(command)

			// print error logs if the response is invalid
			if !success {

				// open the file in append mode; if the file doesn't exist, it will be created
				file, err := os.OpenFile("error_logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				// create the error log string to be added
				errorLogStartTimeBuffer := time.Now().String()
				errorLogGoBuffer := command
				errorLogResponseBuffer := response

				errorString := errorLogStartTimeBuffer + ": ERROR IN SEARCH. SELECTED A RANDOM MOVE. PREVIOUS COMMANDS RECEIVED:\n" + errorLogPositionBuffer + "\n" + errorLogGoBuffer + "\n" + errorLogResponseBuffer + "\n"

				// now add the log strings to the error log file
				_, err = fmt.Fprintln(file, errorString)
				if err != nil {
					log.Fatal(err)
				}
			}

			// send the best move response
			fmt.Printf("%v\n", response)

		} else if strings.HasPrefix(command, "stop") {
			pos.command_stop()

		} else if strings.HasPrefix(command, "ponderhit") {
			pos.command_ponderHit()

		} else if strings.HasPrefix(command, "quit") {
			pos.command_quit()
			runCommandLoop = false
		}
	}
}

// --------------------------------------------------------- UCI -----------------------------------------------
/*
GUI to engine:
--------------

  - uci
    tell engine to use the uci (universal chess interface),
    this will be sent once as a first command after program boot
    to tell the engine to switch to uci mode.
    After receiving the uci command the engine must identify itself with the "id" command
    and send the "option" commands to tell the GUI which engine settings the engine supports if any.
    After that the engine should send "uciok" to acknowledge the uci mode.
    If no uciok is sent within a certain time period, the engine task will be killed by the GUI.

Engine to GUI:
--------------
* id

  - name <x>
    this must be sent after receiving the "uci" command to identify the engine,
    e.g. "id name Shredder X.Y\n"

  - author <x>
    this must be sent after receiving the "uci" command to identify the engine,
    e.g. "id author Stefan MK\n"

Engine to GUI:
--------------
* uciok
	Must be sent after the id and optional options to tell the GUI that the engine
	has sent all infos and is ready in uci mode.


Engine to GUI:
--------------
* option
	This command tells the GUI which parameters can be changed in the engine.
	This should be sent once at engine startup after the "uci" and the "id" commands
	if any parameter can be changed in the engine.
	The GUI should parse this and build a dialog for the user to change the settings.
	Note that not every option needs to appear in this dialog as some options like
	"Ponder", "UCI_AnalyseMode", etc. are better handled elsewhere or are set automatically.
	If the user wants to change some settings, the GUI will send a "setoption" command to the engine.
	Note that the GUI need not send the setoption command when starting the engine for every option if
	it doesn't want to change the default value.
	For all allowed combinations see the examples below,
	as some combinations of this tokens don't make sense.
	One string will be sent for each parameter.
	* name <id>
		The option has the name id.
		Certain options have a fixed value for <id>, which means that the semantics of this option is fixed.
		Usually those options should not be displayed in the normal engine options window of the GUI but
		get a special treatment. "Pondering" for example should be set automatically when pondering is
		enabled or disabled in the GUI options. The same for "UCI_AnalyseMode" which should also be set
		automatically by the GUI. All those certain options have the prefix "UCI_" except for the
		first 6 options below. If the GUI gets an unknown Option with the prefix "UCI_", it should just
		ignore it and not display it in the engine's options dialog.
		* <id> = Hash, type is spin
			the value in MB for memory for hash tables can be changed,
			this should be answered with the first "setoptions" command at program boot
			if the engine has sent the appropriate "option name Hash" command,
			which should be supported by all engines!
			So the engine should use a very small hash first as default.
		* <id> = NalimovPath, type string
			this is the path on the hard disk to the Nalimov compressed format.
			Multiple directories can be concatenated with ";"
		* <id> = NalimovCache, type spin
			this is the size in MB for the cache for the nalimov table bases
			These last two options should also be present in the initial options exchange dialog
			when the engine is booted if the engine supports it
		* <id> = Ponder, type check
			this means that the engine is able to ponder.
			The GUI will send this whenever pondering is possible or not.
			Note: The engine should not start pondering on its own if this is enabled, this option is only
			needed because the engine might change its time management algorithm when pondering is allowed.
		* <id> = OwnBook, type check
			this means that the engine has its own book which is accessed by the engine itself.
			if this is set, the engine takes care of the opening book and the GUI will never
			execute a move out of its book for the engine. If this is set to false by the GUI,
			the engine should not access its own book.
		* <id> = MultiPV, type spin
			the engine supports multi best line or k-best mode. the default value is 1
		* <id> = UCI_ShowCurrLine, type check, should be false by default,
			the engine can show the current line it is calculating. see "info currline" above.
		* <id> = UCI_ShowRefutations, type check, should be false by default,
			the engine can show a move and its refutation in a line. see "info refutations" above.
		* <id> = UCI_LimitStrength, type check, should be false by default,
			The engine is able to limit its strength to a specific Elo number,
		   This should always be implemented together with "UCI_Elo".
		* <id> = UCI_Elo, type spin
			The engine can limit its strength in Elo within this interval.
			If UCI_LimitStrength is set to false, this value should be ignored.
			If UCI_LimitStrength is set to true, the engine should play with this specific strength.
		   This should always be implemented together with "UCI_LimitStrength".
		* <id> = UCI_AnalyseMode, type check
		   The engine wants to behave differently when analysing or playing a game.
		   For example when playing it can use some kind of learning.
		   This is set to false if the engine is playing a game, otherwise it is true.
		 * <id> = UCI_Opponent, type string
		   With this command the GUI can send the name, title, elo and if the engine is playing a human
		   or computer to the engine.
		   The format of the string has to be [GM|IM|FM|WGM|WIM|none] [<elo>|none] [computer|human] <name>
		   Examples:
		   "setoption name UCI_Opponent value GM 2800 human Gary Kasparov"
		   "setoption name UCI_Opponent value none none computer Shredder"
		 * <id> = UCI_EngineAbout, type string
		   With this command, the engine tells the GUI information about itself, for example a license text,
		   usually it doesn't make sense that the GUI changes this text with the setoption command.
		   Example:
			"option name UCI_EngineAbout type string default Shredder by Stefan Meyer-Kahlen, see www.shredderchess.com"
		* <id> = UCI_ShredderbasesPath, type string
			this is either the path to the folder on the hard disk containing the Shredder endgame databases or
			the path and filename of one Shredder endgame datbase.
	   * <id> = UCI_SetPositionValue, type string
	      the GUI can send this to the engine to tell the engine to use a certain value in centipawns from white's
	      point of view if evaluating this specifix position.
	      The string can have the formats:
	      <value> + <fen> | clear + <fen> | clearall

	* type <t>
		The option has type t.
		There are 5 different types of options the engine can send
		* check
			a checkbox that can either be true or false
		* spin
			a spin wheel that can be an integer in a certain range
		* combo
			a combo box that can have different predefined strings as a value
		* button
			a button that can be pressed to send a command to the engine
		* string
			a text field that has a string as a value,
			an empty string has the value "<empty>"
	* default <x>
		the default value of this parameter is x
	* min <x>
		the minimum value of this parameter is x
	* max <x>
		the maximum value of this parameter is x
	* var <x>
		a predefined value of this parameter is x
	Examples:
    Here are 5 strings for each of the 5 possible types of options
	   "option name Nullmove type check default true\n"
      "option name Selectivity type spin default 2 min 0 max 4\n"
	   "option name Style type combo default Normal var Solid var Normal var Risky\n"
	   "option name NalimovPath type string default c:\\n"
	   "option name Clear Hash type button\n"

Engine to GUI:
--------------
* copyprotection
	this is needed for copyprotected engines. After the uciok command the engine can tell the GUI,
	that it will check the copy protection now. This is done by "copyprotection checking".
	If the check is ok the engine should send "copyprotection ok", otherwise "copyprotection error".
	If there is an error the engine should not function properly but should not quit alone.
	If the engine reports "copyprotection error" the GUI should not use this engine
	and display an error message instead!
	The code in the engine can look like this
      TellGUI("copyprotection checking\n");
	   // ... check the copy protection here ...
	   if(ok)
	      TellGUI("copyprotection ok\n");
      else
         TellGUI("copyprotection error\n");

*/

func (pos *Position) command_uci() {
	// at this point we do not initialize support tables yet
	// just acknowledge uci as the communication protocol

	// <<< 1 >>> ID the engine
	fmt.Printf("id name Invincible 1.0\n")
	fmt.Printf("id author Johan vD\n")

	// <<< 2 >>> Options
	// none for now

	// <<< 2 >>> Final response
	fmt.Printf("uciok\n")

}

// --------------------------------------------------------- Debug -----------------------------------------------
/*
GUI to engine:
--------------
  - debug [ on | off ]
    switch the debug mode of the engine on and off.
    In debug mode the engine should send additional infos to the GUI, e.g. with the "info string" command,
    to help debugging, e.g. the commands that the engine has received etc.
    This mode should be switched off by default and this command can be sent
    any time, also when the engine is thinking.
*/
func (pos *Position) command_debug() {
	// ignore for now
}

// --------------------------------------------------------- Is Ready -----------------------------------------------
/*
GUI to engine:
--------------
  - isready
    this is used to synchronize the engine with the GUI. When the GUI has sent a command or
    multiple commands that can take some time to complete,
    this command can be used to wait for the engine to be ready again or
    to ping the engine to find out if it is still alive.
    E.g. this should be sent after setting the path to the tablebases as this can take some time.
    This command is also required once before the engine is asked to do any search
    to wait for the engine to finish initializing.
    This command must always be answered with "readyok" and can be sent also when the engine is calculating
    in which case the engine should also immediately answer with "readyok" without stopping the search.

Engine to GUI
--------------
  - readyok
    This must be sent when the engine has received an "isready" command and has
    processed all input and is ready to accept new commands now.
    It is usually sent after a command that can take some time to be able to wait for the engine,
    but it can be used anytime, even when the engine is searching,
    and must always be answered with "isready".
*/
func (pos *Position) command_isReady() {
	fmt.Printf("readyok\n")
}

// --------------------------------------------------------- Set Option -----------------------------------------------
/*
GUI to engine:
--------------
  - setoption name <id> [value <x>]
    this is sent to the engine when the user wants to change the internal parameters
    of the engine. For the "button" type no value is needed.
    One string will be sent for each parameter and this will only be sent when the engine is waiting.
    The name and value of the option in <id> should not be case sensitive and can inlude spaces.
    The substrings "value" and "name" should be avoided in <id> and <x> to allow unambiguous parsing,
    for example do not use <name> = "draw value".
    Here are some strings for the example below:
    "setoption name Nullmove value true\n"
    "setoption name Selectivity value 3\n"
    "setoption name Style value Risky\n"
    "setoption name Clear Hash\n"
    "setoption name NalimovPath value c:\chess\tb\4;c:\chess\tb\5\n"
*/
func (pos *Position) command_setOption() {
	// no options for now
}

// --------------------------------------------------------- Register -----------------------------------------------
/*
GUI to engine:
--------------
  - register
    this is the command to try to register an engine or to tell the engine that registration
    will be done later. This command should always be sent if the engine	has sent "registration error"
    at program startup.
    The following tokens are allowed:
  - later
    the user doesn't want to register the engine now.
  - name <x>
    the engine should be registered with the name <x>
  - code <y>
    the engine should be registered with the code <y>
    Example:
    "register later"
    "register name Stefan MK code 4359874324"

Engine to GUI:
--------------
  - registration
    this is needed for engines that need a username and/or a code to function with all features.
    Analog to the "copyprotection" command the engine can send "registration checking"
    after the uciok command followed by either "registration ok" or "registration error".
    Also after every attempt to register the engine it should answer with "registration checking"
    and then either "registration ok" or "registration error".
    In contrast to the "copyprotection" command, the GUI can use the engine after the engine has
    reported an error, but should inform the user that the engine is not properly registered
    and might not use all its features.
    In addition the GUI should offer to open a dialog to
    enable registration of the engine. To try to register an engine the GUI can send
    the "register" command.
    The GUI has to always answer with the "register" command	if the engine sends "registration error"
    at engine startup (this can also be done with "register later")
    and tell the user somehow that the engine is not registered.
    This way the engine knows that the GUI can deal with the registration procedure and the user
    will be informed that the engine is not properly registered.
*/
func (pos *Position) command_register() {
	// not needed
}

// --------------------------------------------------------- UCI New Game -----------------------------------------------
/*
GUI to engine:
--------------
  - ucinewgame
    this is sent to the engine when the next search (started with "position" and "go") will be from
    a different game. This can be a new game the engine should play or a new game it should analyse but
    also the next position from a testsuite with positions only.
    If the GUI hasn't sent a "ucinewgame" before the first "position" command, the engine shouldn't
    expect any further ucinewgame commands as the GUI is probably not supporting the ucinewgame command.
    So the engine should not rely on this command even though all new GUIs should support it.
    As the engine's reaction to "ucinewgame" can take some time the GUI should always send "isready"
    after "ucinewgame" to wait for the engine to finish its operation.
*/
func (pos *Position) command_uciNewGame() {
	initEngine()
	pos.reset()

	// "isready" will be sent by the GUI after this command
	// so once the init is done that command will be processed correctly
	// nothing else is needed here

}

// --------------------------------------------------------- Position -----------------------------------------------
/*
GUI to engine:
--------------
  - position [fen <fenstring> | startpos ]  moves <move1> .... <movei>
    set up the position described in fenstring on the internal board and
    play the moves on the internal chess board.
    if the game was played  from the start position the string "startpos" will be sent
    Note: no "new" command is needed. However, if this position is from a different game than
    the last position sent to the engine, the GUI should have sent a "ucinewgame" inbetween.
*/
func (pos *Position) command_position(command string) {
	initEngine()
	pos.reset()

	parts := strings.Split(command, " ")

	var fen string
	for i, part := range parts {

		// parse the fen string
		// this can either be an actual string or the term "startpos"
		if part == "fen" {
			fen = strings.Join(parts[i+1:i+7], " ")
			pos.step1InitFen(fen)
			pos.step2InitRest()
		}

		if part == "startpos" {
			fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			pos.step1InitFen(fen)
			pos.step2InitRest()
		}

		// parse the remaining moves (illegal moves and text will just be ignored)
		pos.makeUCIMove(part)

	}
}

// translates the uci move input to a move recognized by the engine
// it will do nothing if a move is not recognized
func (pos *Position) makeUCIMove(input string) {

	// split the input to the separate string parts
	var fromStr string
	var toStr string
	var promoteStr string

	if len(input) == 4 { // normal moves, including castle moves
		fromStr = input[0:2]
		toStr = input[2:]
	} else if len(input) == 5 { // promotion moves
		fromStr = input[0:2]
		toStr = input[2:4]
		promoteStr = input[4:]
	} else { // unrecognized move
		return
	}

	// convert the string parts to engine ints
	fromSq := getSqFromString(fromStr)
	toSq := getSqFromString(toStr)

	var promoteType int
	if len(promoteStr) == 1 {
		switch promoteStr {
		case "q":
			promoteType = PROMOTION_QUEEN
		case "r":
			promoteType = PROMOTION_ROOK
		case "n":
			promoteType = PROMOTION_KNIGHT
		case "b":
			promoteType = PROMOTION_BISHOP
		}
	}

	// generate moves before a move can be processed
	pos.generateLegalMoves(false)

	// if there are no available moves, return
	//if pos.availableMovesCounter <= 0 {
	if pos.totalMovesCounter <= 0 {
		return
	}

	// else copy the moves
	//allMoves := make([]Move, pos.availableMovesCounter)
	//copy(allMoves, pos.availableMoves[:pos.availableMovesCounter])
	allMoves := make([]Move, pos.totalMovesCounter)
	copy(allMoves, pos.threatMoves[:pos.threatMovesCounter])
	copy(allMoves[pos.threatMovesCounter:], pos.quietMoves[:pos.quietMovesCounter])

	// loop over moves
	// and where the input matches the move, play that move
	var playedMove Move
	foundMove := false
	for _, move := range allMoves {
		if move.getFromSq() == fromSq && move.getToSq() == toSq && move.getPromotionType() == promoteType {
			playedMove = move
			foundMove = true
		}
	}

	if foundMove {
		pos.makeMove(playedMove)
	}
}

// --------------------------------------------------------- Go -----------------------------------------------
/*
GUI to engine:
--------------
  - go
    start calculating on the current position set up with the "position" command.
    There are a number of commands that can follow this command, all will be sent in the same string.
    If one command is not sent its value should be interpreted as it would not influence the search.
  - searchmoves <move1> .... <movei>
    restrict search to this moves only
    Example: After "position startpos" and "go infinite searchmoves e2e4 d2d4"
    the engine should only search the two moves e2e4 and d2d4 in the initial position.
  - ponder
    start searching in pondering mode.
    Do not exit the search in ponder mode, even if it's mate!
    This means that the last move sent in in the position string is the ponder move.
    The engine can do what it wants to do, but after a "ponderhit" command
    it should execute the suggested move to ponder on. This means that the ponder move sent by
    the GUI can be interpreted as a recommendation about which move to ponder. However, if the
    engine decides to ponder on a different move, it should not display any mainlines as they are
    likely to be misinterpreted by the GUI because the GUI expects the engine to ponder
    on the suggested move.
  - wtime <x>
    white has x msec left on the clock
  - btime <x>
    black has x msec left on the clock
  - winc <x>
    white increment per move in mseconds if x > 0
  - binc <x>
    black increment per move in mseconds if x > 0
  - movestogo <x>
    there are x moves to the next time control,
    this will only be sent if x > 0,
    if you don't get this and get the wtime and btime it's sudden death
  - depth <x>
    search x plies only.
  - nodes <x>
    search x nodes only,
  - mate <x>
    search for a mate in x moves
  - movetime <x>
    search exactly x mseconds
  - infinite
    search until the "stop" command. Do not exit the search without being told so in this mode!

Engine to GUI:
--------------
* bestmove <move1> [ ponder <move2> ]
	the engine has stopped searching and found the move <move> best in this position.
	the engine can send the move it likes to ponder on. The engine must not start pondering automatically.
	this command must always be sent if the engine stops searching, also in pondering mode if there is a
	"stop" command, so for every "go" command a "bestmove" command is needed!
	Directly before that the engine should send a final "info" command with the final search information,
	the the GUI has the complete statistics about the last search.
*/
func (pos *Position) command_go(command string) (string, bool) {

	// split the command into a slice of strings
	parts := strings.Split(command, " ")

	// reset the time management variables
	allowedTimeMs := 0
	hasIncrement := false
	allocateTime := true

	// loop over the command parts
	for i, part := range parts {

		// add time to allowed time
		if part == "wtime" {
			if pos.isWhiteTurn {
				timeGain, _ := strconv.Atoi(parts[i+1])
				allowedTimeMs += timeGain
			}
		}

		// add time to allowed time
		if part == "btime" {
			if !pos.isWhiteTurn {
				timeGain, _ := strconv.Atoi(parts[i+1])
				allowedTimeMs += timeGain
			}
		}

		// flag so we know there is increment or not
		if part == "winc" {
			if pos.isWhiteTurn {
				hasIncrement = true
			}
		}

		// flag so we know there is increment or not
		if part == "binc" {
			if !pos.isWhiteTurn {
				hasIncrement = true
			}
		}

		// set the time to this specifically
		if part == "movetime" {
			specificTime, _ := strconv.Atoi(parts[i+1])
			allowedTimeMs = specificTime
			allocateTime = false
		}

		// search infinitely
		// note: the engine still needs to properly implement stop, so this command will likely give errors for now
		if part == "infinite" {
			allowedTimeMs = 1000000 // 1000 sec
			allocateTime = false
		}
	}

	// after getting total time, we allocate time for a search
	// we estimate the game stage using the stage from the position
	// we can take a bit longer in the endgame (normally less moves to play because of less pieces but also less time remaining)
	var timeFactor int
	if hasIncrement {
		if pos.evalMidVsEndStage >= 20 { // opening
			timeFactor = 30
		} else if pos.evalMidVsEndStage >= 16 { // opening-middlegame
			timeFactor = 22
		} else if pos.evalMidVsEndStage >= 12 { // middlegame
			timeFactor = 16
		} else if pos.evalMidVsEndStage >= 8 { // middlegame-endgame
			timeFactor = 12
		} else if pos.evalMidVsEndStage >= 4 { // endgame
			timeFactor = 10
		} else { // endgame try checkmate
			timeFactor = 8
		}
	} else { // if there is no increment, we just set a constant factor
		timeFactor = 24
	}

	// calculate the time we have for the search
	var timeForSearch int
	if allocateTime {
		timeForSearch = allowedTimeMs / timeFactor
	} else {
		timeForSearch = allowedTimeMs
	}

	// we then do the search with the calculated time
	pos.searchForBestMove(timeForSearch)

	// we then get the best move from the search
	success := false
	var bestMove Move
	if pos.bestMove == BLANK_MOVE { // if no move was found, use a random move
		// the issue is that a 3-fold repetition might not be automatically claimed by an opponent in the gui (we assume we will always claim it),
		// and we then need to play on until it is claimed,
		// however, in those cases we never even get to iterate over moves, because we return a search score of 0 early at the root for 3-fold repetitions
		// therefore, in the rare case where 3-fold repetition is not claimed, we still need to add code to manage that
		// for now, we just take the 1st legal move as the best move (we assume it will be rare)
		pos.generateLegalMoves(false)
		if pos.quietMovesCounter > 0 {
			bestMove = pos.quietMoves[0]
		} else {
			bestMove = pos.threatMoves[0]
		}
	} else { // else use the best move found in the search
		bestMove = pos.bestMove
		success = true
	}

	// convert the best move to a format in uci and return it, and the success flag
	moveFromStr := getStringFromSq(bestMove.getFromSq())
	moveToStr := getStringFromSq(bestMove.getToSq())
	promoteStr := getPromotionStringFromType(bestMove.getPromotionType())

	output := "bestmove " + moveFromStr + moveToStr + promoteStr
	return output, success
}

// --------------------------------------------------------- Stop -----------------------------------------------
/*
GUI to engine:
--------------
  - stop
    stop calculating as soon as possible,
    don't forget the "bestmove" and possibly the "ponder" token when finishing the search

Engine to GUI:
--------------
* bestmove <move1> [ ponder <move2> ]
	the engine has stopped searching and found the move <move> best in this position.
	the engine can send the move it likes to ponder on. The engine must not start pondering automatically.
	this command must always be sent if the engine stops searching, also in pondering mode if there is a
	"stop" command, so for every "go" command a "bestmove" command is needed!
	Directly before that the engine should send a final "info" command with the final search information,
	the the GUI has the complete statistics about the last search.

Engine to GUI:
--------------
* info
	the engine wants to send information to the GUI. This should be done whenever one of the info has changed.
	The engine can send only selected infos or multiple infos with one info command,
	e.g. "info currmove e2e4 currmovenumber 1" or
	     "info depth 12 nodes 123456 nps 100000".
	Also all infos belonging to the pv should be sent together
	e.g. "info depth 2 score cp 214 time 1242 nodes 2124 nps 34928 pv e2e4 e7e5 g1f3"
	I suggest to start sending "currmove", "currmovenumber", "currline" and "refutation" only after one second
	to avoid too much traffic.
	Additional info:
	* depth <x>
		search depth in plies
	* seldepth <x>
		selective search depth in plies,
		if the engine sends seldepth there must also be a "depth" present in the same string.
	* time <x>
		the time searched in ms, this should be sent together with the pv.
	* nodes <x>
		x nodes searched, the engine should send this info regularly
	* pv <move1> ... <movei>
		the best line found
	* multipv <num>
		this for the multi pv mode.
		for the best move/pv add "multipv 1" in the string when you send the pv.
		in k-best mode always send all k variants in k strings together.
	* score
		* cp <x>
			the score from the engine's point of view in centipawns.
		* mate <y>
			mate in y moves, not plies.
			If the engine is getting mated use negative values for y.
		* lowerbound
	      the score is just a lower bound.
		* upperbound
		   the score is just an upper bound.
	* currmove <move>
		currently searching this move
	* currmovenumber <x>
		currently searching move number x, for the first move x should be 1 not 0.
	* hashfull <x>
		the hash is x permill full, the engine should send this info regularly
	* nps <x>
		x nodes per second searched, the engine should send this info regularly
	* tbhits <x>
		x positions where found in the endgame table bases
	* sbhits <x>
		x positions where found in the shredder endgame databases
	* cpuload <x>
		the cpu usage of the engine is x permill.
	* string <str>
		any string str which will be displayed be the engine,
		if there is a string command the rest of the line will be interpreted as <str>.
	* refutation <move1> <move2> ... <movei>
	   move <move1> is refuted by the line <move2> ... <movei>, i can be any number >= 1.
	   Example: after move d1h5 is searched, the engine can send
	   "info refutation d1h5 g6h5"
	   if g6h5 is the best answer after d1h5 or if g6h5 refutes the move d1h5.
	   if there is no refutation for d1h5 found, the engine should just send
	   "info refutation d1h5"
		The engine should only send this if the option "UCI_ShowRefutations" is set to true.
	* currline <cpunr> <move1> ... <movei>
	   this is the current line the engine is calculating. <cpunr> is the number of the cpu if
	   the engine is running on more than one cpu. <cpunr> = 1,2,3....
	   if the engine is just using one cpu, <cpunr> can be omitted.
	   If <cpunr> is greater than 1, always send all k lines in k strings together.
		The engine should only send this if the option "UCI_ShowCurrLine" is set to true.
*/
func (pos *Position) command_stop() {
	// nothing extra needed, we already set a stop flag for searches, and return a "bestmove"
}

// --------------------------------------------------------- Ponder Hit -----------------------------------------------
/*
GUI to engine:
--------------
  - ponderhit
    the user has played the expected move. This will be sent if the engine was told to ponder on the same move
    the user has played. The engine should continue searching but switch from pondering to normal search.
*/
func (pos *Position) command_ponderHit() {
	// this engine cannot ponder
}

// --------------------------------------------------------- Quit -----------------------------------------------
/*
GUI to engine:
--------------
  - quit
    quit the program as soon as possible
*/
func (pos *Position) command_quit() {
	// nothing extra needed
}
