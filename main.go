package main

import (
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

type Potential_Client struct {
	Name string
	Email string
	Tel string
	Interest []string
	Message string
}

type Image_File struct {
	Name string
	Width int
	Height int
}

var templates *template.Template

var funcMap = template.FuncMap{
	"noescape": noescape,
	"add": add,
	"minus": minus,
}

func init() {
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
}

func main() {
	//resurce end points
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./content"))))
	http.HandleFunc("/favicon.ico", ServeFileHandler)

	//page end points
	http.HandleFunc("/", ServePublicIndexPage)
	http.HandleFunc("/installments", ServeClientInstallmentsPage)
	http.HandleFunc("/gallery", ServeGalleryPage)
	http.HandleFunc("/team", ServeOurTeamPage)
	http.HandleFunc("/renewable/energy", ServeRenewableEnergyPage)
	http.HandleFunc("/estimates", ServePublicIndexPage)
	http.HandleFunc("/services", ServeOurServicesPage)
	http.HandleFunc("/incentives/financing", ServeOurIncentivesFinanacingPage)
	http.HandleFunc("/tos", ServeOurTOSPage)
	http.HandleFunc("/privacy", ServeOurPrivacyPage)
	http.HandleFunc("/licenses", ServeOurLicensesPage)
	http.HandleFunc("/sitemap", ServeOurSitemapPage)

	//go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	log.Println(http.ListenAndServe(":80", nil))

	log.Println(http.ListenAndServeTLS(":443",
	"/etc/letsencrypt/live/myhytechenergy.com/fullchain.pem",
	"/etc/letsencrypt/live/myhytechenergy.com/privkey.pem",
	nil))
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func add(a int, b int) int {
	return a + b
}

func minus(b int, a int) int {
	return b - a
}

func ServeFileHandler(w http.ResponseWriter, r *http.Request) {
	fname := path.Base(r.URL.Path)
	http.ServeFile(w, r, "./content/static/"+ fname)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		http.StatusTemporaryRedirect)
}

func ServePublicIndexPage(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	pclient := Potential_Client{
		Name: r.FormValue("name"),
		Email: r.FormValue("email"),
		Tel: r.FormValue("tel"),
		Interest: r.Form["interest"],
		Message: r.FormValue("message"),
	}

	msg_flag := ""
	barlast := "0%"
	progress := "25%"
	//pclient := Potential_Client{}

	if r.Method == "POST" {
		if r.FormValue("submit") == "Yes, Send!" {
			Mailer(pclient)
			msg_flag = "finnish"
			barlast = "75%"
			progress = "100%"
		} else if r.FormValue("submit") == "No Thanks" {
			Mailer(pclient)
			msg_flag = "nothanks"
			barlast = "75%"
			progress = "100%"
		} else if r.FormValue("submit") == "Send My Interests" {
			//pclient = Mailer(w, r)
			msg_flag = "message"
			barlast = "50%"
			progress = "75%"
		} else {
			//pclient = Mailer(w, r)
			msg_flag = "interest"
			barlast = "25%"
			progress = "50%"
		}
	}

	fmt.Println(r.FormValue("submit"))
	//l, _ := FileList("imgs")

	//msg_flag = "message"

	rays := make([]int, 20)
	lessRays := make([]int, 10)
	chars := []rune("Imagine an energy bill paid for by the sun");
	var comp []string
	ray_deg := 0

	for i, _ := range rays {
		rays[i] = ray_deg
		ray_deg += 18
	}

	ray_deg = 0
	for i, _ := range lessRays {
		lessRays[i] = ray_deg
		ray_deg += 36
	}

	str := ""

	for i, _ := range chars {
		str = string(chars[:i+1])
		comp = append(comp, str)
	}

	data := map[string]interface{}{
		"Ask": msg_flag,
		"Pclient": pclient,
		"Interests": strings.Join(pclient.Interest, ", "),
		//"Files": l,
		"Page": "",
		"Barlast": barlast,
		"Progress": progress,
		"Rays": rays,
		"LessRays": lessRays,
		"TextRays": comp,
	}


	var err error

	if r.URL.Path == "/estimates" {
		data["Page"] = "/estimates"
		err = templates.ExecuteTemplate(w, "estimates.html", data)

	} else {
		data["Page"] = "/"
		err = templates.ExecuteTemplate(w, "public_index.html", data)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
}

func ServeClientInstallmentsPage(w http.ResponseWriter, r *http.Request) {

	list, err := FileList("imgs")

	if err != nil {
		log.Fatal(err)
		return
	}

	data := map[string]interface{}{
		"Images": list,
		"Count": len(list),
		"Page": "/installments",
	}
	err = templates.ExecuteTemplate(w, "client-installments.html", data)

	if err != nil {
		log.Fatal(err)
		return
	}
}
func ServeGalleryPage(w http.ResponseWriter, r *http.Request) {

	list, err := FileList("imgs")

	if err != nil {
		log.Fatal(err)
		return
	}

	data := map[string]interface{}{
		"Images": list,
		"Page": "/gallery",
	}
	err = templates.ExecuteTemplate(w, "gallery.html", data)

	if err != nil {
		log.Fatal(err)
		return
	}
}
func ServeOurTeamPage(w http.ResponseWriter, r *http.Request) {

	profiles, err := FileList("profiles")

	if err != nil {
		log.Fatal(err)
		return
	}



	data := map[string]interface{}{
		"Profiles": profiles,
		"Page": "/team",
	}
	err = templates.ExecuteTemplate(w, "team.html", data)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func ServeOurServicesPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/services",
	}
	err := templates.ExecuteTemplate(w, "services.html", data)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func FileList(folder string) (list []Image_File, err error) {
	//file, err := os.Open("./content/" + folder)
	//
	//if err != nil {
	//	log.Fatalf("failed opening directory: %s", err)
	//}
	//defer file.Close()
	//
	//return file.Readdirnames(0) // 0 to read all files and folders

	dir_to_scan := "./content/" + folder

	files, _ := ioutil.ReadDir(dir_to_scan)
	for _, imgFile := range files {
		if reader, err := os.Open(filepath.Join(dir_to_scan, imgFile.Name())); err == nil {
			defer reader.Close()
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}
			//fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
			im.Width = im.Width / 2;
			im.Height = im.Height / 2;

			list = append(list, Image_File{imgFile.Name(), im.Width, im.Height})
		} else {
			fmt.Println("Impossible to open the file:", err)
		}
	}
	return list, err
}

func Mailer(p Potential_Client) {

	conf_msg := []byte("Subject: From HyTech Energy" + "\r\n" +
				"\rThank you " + p.Name + " for reaching out to HyTech Energy!\r\n" +
				"\rWe will be contacting you soon about an incredible estimate on your renewable energy future.\r\n")

	msg := []byte("Subject: From " + p.Name + "\r\n" +
		"\rClient: " + p.Name + " <" + p.Email + ">\r\n" +
		"\rEmail: "+ p.Email +"\r\n" +
		"\rPhone: "+ p.Tel +"\r\n" +
		"\rInterests: " + strings.Join(p.Interest, ", ") + "\r\n" +
		"\r\n" + p.Message + "\r\n")

		from := "patchcoding@gmail.com"
		pass := "Field4d in con!"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{"patchkientz@gmail.com"}, []byte(msg))

//lhoehn@hytechnrg.com

	if err != nil {
		log.Printf("smtp error: %s", err)
	}

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{p.Email}, []byte(conf_msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
	}
}

func ServeRenewableEnergyPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/renewable/energy",
	}

	err := templates.ExecuteTemplate(w, "renewable-energy.html", data)

	if err != nil {
		log.Fatal(err)
	}
}

func ServeFreeEstimatesPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "",
	}

	err := templates.ExecuteTemplate(w, "estimates.html", data)

	if err != nil {
		log.Fatal(err)
	}
}
func ServeOurIncentivesFinanacingPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/incentives/financing",
	}

	err := templates.ExecuteTemplate(w, "incentives-financing.html", data)

	if err != nil {
		log.Fatal(err)
	}
}
func ServeOurTOSPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/tos",
	}

	err := templates.ExecuteTemplate(w, "tos.html", data)

	if err != nil {
		log.Fatal(err)
	}
}
func ServeOurPrivacyPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/privacy",
	}

	err := templates.ExecuteTemplate(w, "privacy.html", data)

	if err != nil {
		log.Fatal(err)
	}
}
func ServeOurLicensesPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/licenses",
	}

	err := templates.ExecuteTemplate(w, "licenses.html", data)

	if err != nil {
		log.Fatal(err)
	}
}
func ServeOurSitemapPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"Page": "/sitemap",
	}

	err := templates.ExecuteTemplate(w, "sitemap.html", data)

	if err != nil {
		log.Fatal(err)
	}
}


