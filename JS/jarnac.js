const readline = require('readline');

//le sac contenant les lettres
let sac = [
    ...Array(14).fill('A'),
    ...Array(4).fill('B'),
    ...Array(7).fill('C'),
    ...Array(5).fill('D'),
    ...Array(19).fill('E'),
    ...Array(2).fill('F'),
    ...Array(4).fill('G'),
    ...Array(2).fill('H'),
    ...Array(11).fill('I'),
    ...Array(1).fill('J'),
    ...Array(1).fill('K'),
    ...Array(6).fill('L'),
    ...Array(5).fill('M'),
    ...Array(9).fill('N'),
    ...Array(8).fill('O'),
    ...Array(4).fill('P'),
    ...Array(1).fill('Q'),
    ...Array(10).fill('R'),
    ...Array(7).fill('S'),
    ...Array(9).fill('T'),
    ...Array(8).fill('U'),
    ...Array(2).fill('V'),
    ...Array(1).fill('W'),
    ...Array(1).fill('X'),
    ...Array(1).fill('Y'),
    ...Array(2).fill('Z')
];

class Joueur {
    // class d'un joueur
    constructor() {
        this.plateau = [];  
        this.lettres = [];
        this.points = 0;
        this.piocherLettres(6);
    }

    is_valid(mot){
        // verifier si un mot peut être placé sur le plateau avec les lettres disponibles
        mot = mot.toUpperCase();

        for (let i = 0; i < mot.length; i++) {
            if (this.lettres.indexOf(mot[i].toUpperCase()) === -1) {
                return false;
            }
        }
        if (mot.length < 3) {
            return false;
        }
        return true;
    }

    placerMot(mot){
        // placer un mot sur le plateau
        mot = mot.toUpperCase();
        if (this.is_valid(mot)) {
            for (let i = 0; i < mot.length; i++) {
                this.lettres.splice(this.lettres.indexOf(mot[i]), 1);
            }
            this.plateau.push(mot.toUpperCase());
            this.piocherLettres(1);
            this.points += mot.length**2;
            return true;
        }
        else {
            console.log('Vous ne pouvez pas placer ce mot.');
            return false;
        }
    }

    changerLigne(ligne, mot){
        // changer une ligne du plateau
        mot = mot.toUpperCase();
        let ancien_mot = this.plateau[ligne-1].split('');
        let longueur_ancien_mot = ancien_mot.length;
        let lettres_dispos = this.lettres;

        if (ancien_mot.join('') == mot) {
            console.log('Vous ne pouvez pas remplacer ce mot par le même mot.');
            return false;
        }

        for (let i = 0; i<mot.length; i++){
            if (ancien_mot.includes(mot[i])){
                ancien_mot.splice(ancien_mot.indexOf(mot[i]), 1);
            }
            else if (lettres_dispos.includes(mot[i])){
                lettres_dispos.splice(lettres_dispos.indexOf(mot[i]), 1);
            }
            else {
                console.log('Vous ne pouvez pas remplacer ce mot avec les lettres que vous avez.');
                return false;
            }
        }
        if (ancien_mot.length > 0) {
            console.log('Vous ne pouvez pas remplacer ce mot : le mot saisi ne contient pas toutes les lettres du mot initial.');
            return false;
        }
        this.plateau[ligne-1] = mot;
        this.lettres = lettres_dispos;
        this.points = this.points - longueur_ancien_mot**2 + mot.length**2;
        return true;
    }

    piocherLettres(n){
        // piocher n lettres aléatoirement dans le sac
        for (let i = 0; i < n; i++) {
            let index = Math.floor(Math.random() * sac.length);
            this.lettres.push(sac[index]);
            sac.splice(index, 1);
        }
    }

    afficherPlateau(){
        // afficher le plateau
        for (let i = 0; i < this.plateau.length; i++) {
            console.log(String(i+1) + '  ' + this.plateau[i]);
        }
        for (let i = this.plateau.length; i < 8; i++) {
            console.log(String(i+1) + ' ');
        }
    }

    echangerLettres(lettres){
        // échanger 3 de ses lettres avec 3 lettres du sac
        if (lettres.length === 3 && this.is_valid(lettres)) {
            for (let i = 0; i < 3; i++) {
                let index = this.lettres.indexOf(lettres[i].toUpperCase());
                this.lettres.splice(index, 1);
                sac.push(lettres[i].toUpperCase());
            }
            this.piocherLettres(3);
            return true;
        }
        else {
            console.log('Vous ne pouvez pas échanger ces lettres.');
            return false;
        }
    }
}

class Game {
    // class du jeu Jarnac
    constructor() {
        this.joueurs = [];
        this.joueurs.push(new Joueur());
        this.joueurs.push(new Joueur());
    }

    afficherPlateaux(tour) {
        // afficher les plateaux des joueurs 
        console.log('Voici le plateau du joueur 1 :');
        this.joueurs[0].afficherPlateau();
        console.log('Voici vos lettres : ' + this.joueurs[0].lettres.join(' '));
        console.log('Points : ' + this.joueurs[0].points + '\n');

        console.log('Voici le plateau du joueur 2 :');
        this.joueurs[1].afficherPlateau();
        console.log('Voici vos lettres : ' + this.joueurs[1].lettres.join(' '));
        console.log('Points : ' + this.joueurs[1].points + '\n');

        console.log('Joueur ' + (tour % 2 + 1) + ', c\'est à vous de jouer !\n');
    }

    async jarnac(tour) {
        // jarnac : fonction qui permet de voler un mot à l'adversaire
        let rep;
        do {
            let mot;
            let ligne;
            let voleur = this.joueurs[tour % 2];
            let victime = this.joueurs[(tour+1) % 2];
            do {ligne = await this.questionAsync('Quelle ligne voulez-vous voler ? (1-' + victime.plateau.length + ') ')
            } while (ligne < 1 || ligne > victime.plateau.length || isNaN(ligne));

            do {mot = await this.questionAsync('Quel mot voyez-vous ? ')
            } while (!victime.changerLigne(ligne, mot));

            mot = mot.toUpperCase();
            victime.plateau.splice(ligne-1, 1);
            victime.points -= mot.length**2;
            voleur.plateau.push(mot);
            voleur.points += mot.length**2;

            console.clear();
            this.afficherPlateaux(tour);

            do { rep = await this.questionAsync('Voulez-vous voler un autre mot ? (Y/N) ');
            } while (rep !== 'Y' && rep !== 'y' && rep !== 'N' && rep !== 'n');
        } while (rep === 'Y' || rep === 'y');
    }

    async run() {
        // fonction principale du jeu Jarnac
        let fin = false;
        let tour = 0;

        console.clear();
        console.log('Bienvenue dans le jeu Jarnac !\n');

        while (!fin) {
            //boucle principale du jeu : continue jusqu'à ce qu'un des joueurs ait 8 mots sur son plateau
            let joueur = this.joueurs[tour % 2];
            let jarnac = true;
            let tour_de_jeu = true;
            let rep;
            let piocher = true;

            while (tour_de_jeu) {
                // boucle de jeu : continue jusqu'à ce que le joueur passe
                this.afficherPlateaux(tour);
                
                // on peut jarnac au début du tour dès le 2ème tour
                if (jarnac && tour > 0) {
                    do {rep = await this.questionAsync('Voulez-vous Jarnac ? (Y/N) ')
                    } while (rep !== 'Y' && rep !== 'y' && rep !== 'N' && rep !== 'n');
                    if (rep === 'Y' || rep === 'y') {
                        console.log('Jarnac !');
                        await this.jarnac(tour);
                        jarnac = false;
                    }
                    else if (rep === 'N' || rep === 'n') {
                        jarnac = false;
                    }
                }
                
                // on peut choisir de piocher ou d'échanger des lettres dès le 2ème tour
                if (tour > 1 && piocher) {
                    do {rep = await this.questionAsync('Voulez-vous piocher une lettre (P) ou bien échanger 3 de vos lettres avec 3 lettres du sac (E) ? ')
                    } while (rep !== 'P' && rep !== 'p' && rep !== 'E' && rep !== 'e');

                    if (rep === 'P' || rep === 'p') {
                        joueur.piocherLettres(1);
                    }
                    else if (rep === 'E' || rep === 'e') {
                        if (joueur.lettres.length < 3){
                            console.log('Vous n\'avez pas assez de lettres pour échanger.');
                        }
                        else {
                            let lettres_saisies;
                            do {lettres_saisies = await this.questionAsync('Quelles lettres voulez-vous échanger ? (3 lettres) ')
                            } while (!joueur.echangerLettres(lettres_saisies));
                        }
                    }
                    console.clear();
                    this.afficherPlateaux(tour);
                    piocher = false;
                }

                // ici on place un mot, on change une ligne ou on passe
                do {rep = await this.questionAsync('Que voulez-vous faire ? (Placer un mot : M; Changer une ligne : C; Passer : P) ')
                } while (rep !== 'M' && rep !== 'm' && rep !== 'C' && rep !== 'c' && rep !== 'P' && rep !== 'p');

                if (rep === 'M' || rep === 'm') {
                    let mot;
                    do {mot = await this.questionAsync('Quel mot voulez-vous placer ? ')
                    } while (!joueur.placerMot(mot));
                }
                else if (rep === 'C' || rep === 'c') {
                    let ligne;
                    do {ligne = await this.questionAsync('Quelle ligne voulez-vous changer ? (1-' + joueur.plateau.length + ') ')
                    } while (ligne < 1 || ligne > joueur.plateau.length || isNaN(ligne));
                    let mot;
                    do {mot = await this.questionAsync('Quel mot voulez-vous placer ? ')
                    } while (!joueur.changerLigne(ligne, mot));
                }
                else if (rep === 'P' || rep === 'p') {
                    tour_de_jeu = false;
                }
                console.clear();
            }
            if (this.joueurs[0].plateau.length === 8 || this.joueurs[1].plateau.length === 8) {
                // fin de la partie si un des joueurs a 8 mots sur son plateau
                fin = true;
            }
            tour++;
            console.clear();
        }
        console.log('Fin de la partie !\n');
        console.log('Le juoueur 1 a fini avec ' + this.joueurs[0].points + ' points.');
        console.log('Le juoueur 2 a fini avec ' + this.joueurs[1].points + ' points.');
        console.log('Le gagnant est le joueur ' + (this.joueurs[0].points > this.joueurs[1].points ? 1 : 2) + ' !');
    }

    questionAsync(prompt) {
        return new Promise((resolve) => {
            const rl = readline.createInterface({
                input: process.stdin,
                output: process.stdout
            });

            rl.question(prompt, (answer) => {
                rl.close();
                resolve(answer);
            });
        });
    }
}

// lancer le jeu
let game = new Game();
game.run();