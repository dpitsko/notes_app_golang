package main

import (
	"net/http"
    "html/template"
    "log"
    "time"
    "strconv"
    
    "github.com/gorilla/mux"
)


// ------------- DATA STORE ------------
type Note struct {
    Title string
    Description string
    CreatedOn time.Time
}
//View Model for edit
type EditNote struct {
    Note
    Id string
}

//Store for the Notes collection
var noteStore = make(map[string]Note)
//Variable to generate key for the collection
var id int
// --------------------------------------


// ------------- TEMPLATES ---------------
var templates map[string]*template.Template
var baseFile string = "templates/base.goml"

// Compile templates
func init(){
    if templates == nil {
        templates = make(map[string]*template.Template)
    }
    templates["index"] = template.Must(template.ParseFiles("templates/index.goml", baseFile))
    templates["add"] = template.Must(template.ParseFiles("templates/add.goml", baseFile))
    templates["edit"] = template.Must(template.ParseFiles("templates/edit.goml", baseFile))
}

// Render template helper function
func renderTemplate(w http.ResponseWriter, name string, template string, context interface{}) {
    //get template from compiled map holding all templates
    tmpl, ok := templates[name]
    if !ok {
        http.Error(w, "The page does not exist", http.StatusInternalServerError)
    }
    //execute (render) template
    err := tmpl.ExecuteTemplate(w, template, context)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// -------------------------------------

// ------------- MAIN FUNCTION ------------
func main(){

    router := mux.NewRouter().StrictSlash(true)
    
    fs := http.FileServer(http.Dir("static"))  
    router.Handle("/static/", http.StripPrefix("/static/", fs))
    
    router.HandleFunc("/", getNotes)
    router.HandleFunc("/notes/add", addNote)
    router.HandleFunc("/notes/save", saveNote)
    router.HandleFunc("/notes/edit/{id}", editNote)
    router.HandleFunc("/notes/update/{id}", updateNote)
    router.HandleFunc("/notes/delete/{id}", deleteNote)
    
    server := &http.Server{
        Addr: ":8080",
        Handler: router,
    }
    
    log.Println("Listening...")
    server.ListenAndServe()
}
// -------------------------------------

// ------------ HANDLERS ---------------
func getNotes(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index", "base", noteStore)
}

func addNote(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "add", "base", nil)
}

func editNote(w http.ResponseWriter, r *http.Request) {
    var context EditNote
    //get url params
    params := mux.Vars(r)
    key := params["id"]
    
    //get note from map and add it to the view context 
    if note, ok := noteStore[key]; ok {
        context = EditNote{note, key}
    } else {
        http.Error(w, "Note does not exist!", http.StatusBadRequest)
    }
    renderTemplate(w, "edit", "base", context)
}



 //this function does not render a view, it jst handles the form data from /notes/add
func saveNote(w http.ResponseWriter, r *http.Request) {
    
    r.ParseForm()
    title := r.PostFormValue("title")
    desc := r.PostFormValue("description")
    created := time.Now()
    // create note instance based on form data
    note := Note{
        title,
        desc,
        created,
    }
    //increment id and add note to the data store map
    id++
    key := strconv.Itoa(id)
    noteStore[key] = note
    
    //redirect user
    http.Redirect(w, r, "/", 302)
}

 //this function does not render a view, it jst handles the form data from /notes/update/{id}
 func updateNote(w http.ResponseWriter, r *http.Request) {
    
    //get url params
    params := mux.Vars(r)
    key := params["id"]
    var noteToUpdate Note
    
    if _, ok := noteStore[key]; ok { //check if note exists
        
        r.ParseForm()
        //update values from the edit form
        noteToUpdate.Title = r.PostFormValue("title")
        noteToUpdate.Description = r.PostFormValue("description")
        
        //delete existing item and add updated note
        delete(noteStore, key)
        noteStore[key] = noteToUpdate
        
    } else {
        http.Error(w, "Note does not exist!", http.StatusBadRequest)
    }
    
    //redirect user
    http.Redirect(w, r, "/", 302)

 }
 
 //this function does not render a view, it jst handles the form data from /notes/delete/{id}
func deleteNote(w http.ResponseWriter, r *http.Request) {
    //Read value from route variable
    params := mux.Vars(r)
    key := params["id"]
    // Remove from Store
    if _, ok := noteStore[key]; ok {
    //delete existing item
    delete(noteStore, key)
    } else {
        http.Error(w, "Could not find the resource to delete.", http.StatusBadRequest)
    }
    http.Redirect(w, r, "/", 302)
}

// -------------------------------------