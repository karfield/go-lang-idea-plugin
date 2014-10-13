package ro.redeul.google.go.lang.psi.impl;

import com.intellij.codeInsight.lookup.LookupElementBuilder;
import com.intellij.extapi.psi.PsiElementBase;
import com.intellij.lang.ASTNode;
import com.intellij.lang.Language;
import com.intellij.openapi.roots.ex.ProjectRootManagerEx;
import com.intellij.openapi.util.TextRange;
import com.intellij.openapi.vfs.VirtualFile;
import com.intellij.psi.PsiElement;
import com.intellij.psi.PsiFile;
import com.intellij.psi.PsiManager;
import com.intellij.psi.ResolveState;
import com.intellij.psi.scope.PsiScopeProcessor;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.Nullable;
import ro.redeul.google.go.GoLanguage;
import ro.redeul.google.go.lang.psi.GoFile;
import ro.redeul.google.go.lang.psi.GoPackage;
import ro.redeul.google.go.lang.psi.GoPsiElement;
import ro.redeul.google.go.lang.psi.visitors.GoElementVisitor;
import ro.redeul.google.go.lang.psi.visitors.GoElementVisitorWithData;

import java.util.Collection;

public class GoPackageImpl extends PsiElementBase implements GoPackage {

    private PsiManager myPsiManager;
    private final String mySourceRootPath;
    private final String myPath;

    public GoPackageImpl(String importPath, String sourceRootPath, PsiManager psiManager) {
        myPath = importPath;
        mySourceRootPath = sourceRootPath;
        myPsiManager = psiManager;
    }

    @Override
    public void accept(GoElementVisitor visitor) {
        visitor.visitPackage(this);
    }

    @Override
    public <T> T accept(GoElementVisitorWithData<T> visitor) {
        accept((GoElementVisitor) visitor);
        return visitor.getData();
    }

    @Override
    public void acceptChildren(GoElementVisitor visitor) {

    }

    @Override
    public LookupElementBuilder getCompletionPresentation() {
        return null;
    }

    @Override
    public LookupElementBuilder getCompletionPresentation(GoPsiElement child) {
        return null;
    }

    @Override
    public String getPresentationText() {
        return null;
    }

    @Override
    public String getPresentationTailText() {
        return null;
    }

    @Override
    public String getPresentationTypeText() {
        return null;
    }

    @NotNull
    @Override
    public Language getLanguage() {
        return GoLanguage.INSTANCE;
    }

    @NotNull
    @Override
    public PsiElement[] getChildren() {
        return new PsiElement[0];
    }

    @Override
    public PsiElement getParent() {
        return null;
    }

    @Override
    public PsiElement getFirstChild() {
        return null;
    }

    @Override
    public PsiElement getLastChild() {
        return null;
    }

    @Override
    public PsiElement getNextSibling() {
        return null;
    }

    @Override
    public PsiElement getPrevSibling() {
        return null;
    }

    @Override
    public TextRange getTextRange() {
        return null;
    }

    @Override
    public int getStartOffsetInParent() {
        return 0;
    }

    @Override
    public int getTextLength() {
        return 0;
    }

    @Nullable
    @Override
    public PsiElement findElementAt(int offset) {
        return null;
    }

    @Override
    public int getTextOffset() {
        return 0;
    }

    @Override
    public String getText() {
        return null;
    }

    @NotNull
    @Override
    public char[] textToCharArray() {
        return new char[0];
    }

    @Override
    public boolean textContains(char c) {
        return false;
    }

    @Override
    public ASTNode getNode() {
        return null;
    }

    @Override
    @NotNull
    public PsiManager getManager() {
        return myPsiManager;
    }

    @Override
    public String getImportPath() {
        return myPath;
    }

    public String getName() {

        VirtualFile sourceRoots[] = ProjectRootManagerEx.getInstanceEx(getProject()).getContentSourceRoots();

        for (VirtualFile sourceRoot : sourceRoots) {
            VirtualFile packageFolder = sourceRoot.findFileByRelativePath(getImportPath());
            if (packageFolder != null && packageFolder.isDirectory()) {
                VirtualFile files[] = packageFolder.getChildren();
                for (VirtualFile file : files) {
                    PsiFile psiFile = getManager().findFile(file);
                    if (psiFile != null && psiFile instanceof GoFile) {
                        GoFile goFile = (GoFile) psiFile;
                        return goFile.getPackage().getPackageName();
                    }
                }
            }
        }

        return super.getName();
    }

    @Override
    public Collection<GoFile> getPackageFiles() {
        return null;
    }

    @Override
    public boolean processDeclarations(@NotNull PsiScopeProcessor processor,
                                       @NotNull ResolveState state, PsiElement lastParent,
                                       @NotNull PsiElement place) {
        // TODO: implement resolving here
        return true;
    }
}
