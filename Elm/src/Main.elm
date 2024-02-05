module Main exposing (main)

import Array exposing (fromList, get)
import Browser
import Html exposing (Html, div, text, input, button, h1, p, blockquote,ul,li,ol)
import Html.Attributes exposing (placeholder, value, disabled, class, type_, checked)
import Html.Events exposing (onClick, onInput)
import Http exposing (Error)
import Random exposing (int)
import Json.Decode exposing (Decoder, string, list, field, map, map2)
import Debug


-- Initialisation

init : () -> (Model, Cmd Msg)
init _ =
    ( initialModel
    , Http.get {
        url = "../static/mots.txt",
        expect = Http.expectString GotWords
    }
    )

loadJson : String -> Cmd Msg
loadJson word =
    Http.get {
        url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word, 
        expect = Http.expectJson GotJson (list meaningsDecoder)
    }

initialModel : Model
initialModel =
    { selectedWord = Nothing
    , loading = True
    , userInput = ""
    , randomInt = 0
    , wordFound = False
    , words = []
    , checkBoxChecked = False
    , result = Ok []
    }

--Types

type alias Model =
    { selectedWord : Maybe String
    , loading : Bool
    , userInput : String
    , randomInt : Int
    , wordFound : Bool
    , words : List String
    , checkBoxChecked : Bool
    , result : Result Error (List Meanings)
    }

type Msg
    = Noop
    | GotJson (Result Error (List Meanings))
    | GotWords (Result Error String)
    | UserInput String
    | RandomInt Int
    | ToggleCheckBox

type alias Meanings = 
    { word : String
    , meanings : List Meaning
    }
    
type alias Meaning = 
    { partOfSpeech : String
    , definitions : List Definition
    }

type alias Definition =
  { definition : String
  }


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        Noop -> 
            (model, Cmd.none)

        GotWords result ->  
            case result of
                Err reason ->  -- Si on arrive pas à trouver le fichier mots.txt
                    let
                        msg0 = Debug.log "failed to get word" reason
                    in
                    ({model | loading = False}, Cmd.none)
                Ok words ->  -- Si on y arrive
                    let 
                        newWords = String.split " " words
                        randomIndex = Random.generate RandomInt (Random.int 0 (List.length newWords))
                    in
                    ({model | words = newWords}, randomIndex)
        
        RandomInt ri -> -- Ici on choisi un mot aléatoire dans la liste créée au dessus avec l'indice aléatoire obtenu
            let 
                selectedWord = get ri (fromList model.words)
            in
            case selectedWord of
                Just word ->
                    ({model | selectedWord = selectedWord}, loadJson word)
                Nothing ->
                    (model, Cmd.none)
        
        GotJson result ->  -- On vérifie si le json a bien été trouvé et décodé
            case result of
                Err reason ->
                    let
                        msg0 = Debug.log "failed to get the json datas" reason
                    in
                    ({model | loading = False}, Cmd.none)
                Ok json -> 
                    ({model | result = result, loading = False}, Cmd.none)

        UserInput inputWord -> -- Si on modifie le champ pour changer le mot, on vérifie si le mot rentré est le bon
            let
                isCorrect = 
                    case model.selectedWord of
                        Just correctWord ->
                            inputWord == correctWord
                        Nothing ->
                            False
            in
            if isCorrect then
                ({model | userInput = inputWord, wordFound = True}, Cmd.none)
            else
                ({model | userInput = inputWord, wordFound = False}, Cmd.none)
        
        ToggleCheckBox -> -- On gère la checkbox pour montrer le mot ou pas
            ({model | checkBoxChecked = not model.checkBoxChecked}, Cmd.none)

-- View
view : Model -> Html Msg
view model =
    case model.result of
        Ok meanings ->
            if model.loading then div [][text "Loading..."]
            else
                div []
                        [ -- Main title (h1) : changes if word found
                        h1 [] [ text (if model.wordFound then "Well done, bruh !" else "Guess the word, bruh !") ]
                        , -- User input field
                        input [ placeholder "Write the word", value model.userInput, onInput UserInput, class "user-input" ] []
                        , -- Checkbox
                        div [] [text (if model.checkBoxChecked then Maybe.withDefault "No word ......" model.selectedWord else "Show the word :"), input [type_ "checkbox", checked model.checkBoxChecked, onClick ToggleCheckBox][]]
                        , -- Section to display the meanings
                        div []  (List.map viewMeanings meanings)
                        ]
        Err _ ->
            div []
                [ text "Error loading meanings" ]

viewMeanings : Meanings -> Html Msg
viewMeanings meanings =
    div []
        [ 
            div [] (List.map viewMeaning meanings.meanings)
        ]

viewMeaning : Meaning -> Html Msg
viewMeaning meaning =
    div []
        [ ul []
            [ li [] [ text meaning.partOfSpeech ]
            , -- Render the definitions for the meaning
              ol [] (List.map viewDefinition meaning.definitions)
            ]
        ]

viewDefinition : Definition -> Html Msg
viewDefinition definition =
    li []
        [ blockquote []
            [ p [] [ text definition.definition ]]
        ]

-- Decoders

meaningsDecoder : Decoder Meanings
meaningsDecoder =
    map2 Meanings
        (field "word" string)
        (field "meanings" (list meaningDecoder))

meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionDecoder))

definitionDecoder : Decoder Definition
definitionDecoder =
    map Definition
        (field "definition" string)

-- Subscriptions

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none

-- Main

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
